package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var ErrInvalidDuration = errors.New("invalid duration")

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if val, ok := v.(float64); ok {
		*d = Duration(time.Duration(val))

		return nil
	}

	if val, ok := v.(string); ok {
		t, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration(t)

		return nil
	}

	return ErrInvalidDuration
}

type Input struct {
	ContentLength   int64    `json:"content-length"`
	Address         string   `json:"address"`
	XPathExpression string   `json:"xpath-expression"`
	Timeout         Duration `json:"timeout"`
}

type Output struct {
	Success      bool     `json:"success"`
	ErrorMessage string   `json:"error-message,omitempty"`
	Nodes        []string `json:"nodes"`
}

func main() {
	var (
		in  = &Input{Timeout: Duration(time.Second)}
		out = &Output{Success: true}

		err error
	)

	defer func(err *error){
		if *err == nil {
			return
		}

		if encErr := json.NewEncoder(os.Stdout).Encode(&Output{ErrorMessage: (*err).Error()}); encErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, encErr)
		}
	}(&err)

	if err = json.NewDecoder(os.Stdin).Decode(in); err != nil {
		return
	}

	conn, err := net.Dial("tcp", in.Address)
	if err != nil {
		return
	}

	defer conn.Close()

	if err = conn.SetReadDeadline(time.Now().Add(time.Duration(in.Timeout))); err != nil {
		return
	}

	buf := &bytes.Buffer{}
	if _, err = io.CopyN(buf, conn, in.ContentLength); err != nil {
		return
	}

	n, err := htmlquery.Parse(buf)
	if err != nil {
		return
	}

	nn, err := htmlquery.QueryAll(n, in.XPathExpression)
	if err != nil {
		return
	}

	out.Nodes = make([]string, 0, len(nn))
	nbuf := &bytes.Buffer{}

	for _, n := range nn {
		if err = html.Render(nbuf, n); err != nil {
			return
		}

		out.Nodes = append(out.Nodes, html.UnescapeString(nbuf.String()))
		nbuf.Reset()
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	err = enc.Encode(out)
}
