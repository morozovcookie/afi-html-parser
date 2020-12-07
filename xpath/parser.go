package xpath

import (
	"bytes"
	"io"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Parser struct {
	expression string
}

func NewParser(expression string) *Parser {
	return &Parser{
		expression: expression,
	}
}

func (p *Parser) Parse(r io.Reader) ([]string, error) {
	n, err := htmlquery.Parse(r)
	if err != nil {
		return nil, err
	}

	nn, err := htmlquery.QueryAll(n, p.expression)
	if err != nil {
		return nil, err
	}

	var (
		out  = make([]string, 0, len(nn))
		nbuf = &bytes.Buffer{}
	)

	for _, n := range nn {
		if err = html.Render(nbuf, n); err != nil {
			return nil, err
		}

		out = append(out, html.UnescapeString(nbuf.String()))
		nbuf.Reset()
	}

	return out, nil
}
