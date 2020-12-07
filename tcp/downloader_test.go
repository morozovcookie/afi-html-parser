package tcp

import (
	"errors"
	"io"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	ahp "github.com/morozovcookie/afihtmlparser"
)

func TestDownloader_Download(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		dialTimeout time.Duration

		startServer      func(string) (net.Listener, error)
		acceptConnection func(net.Listener) error
		stopServer       func(net.Listener) error

		contentLength int64
		timeout       time.Duration
		callback      ahp.DownloadCallback

		wantErr bool
	}{
		{
			name:    "pass",
			enabled: true,

			dialTimeout: time.Minute,

			startServer: func(address string) (net.Listener, error) {
				ln, err := net.Listen("tcp", address)
				if err != nil {
					return ln, err
				}

				return ln, nil
			},
			acceptConnection: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				conn, err := ln.Accept()
				if err != nil {
					if err == io.EOF {
						return nil
					}

					return err
				}

				defer conn.Close()

				_, err = conn.Write([]byte(`1111111111`))
				if err != nil {
					return err
				}

				return nil
			},
			stopServer: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				return ln.Close()
			},

			contentLength: 10,
			timeout:       time.Second,
			callback: func(r io.Reader) (err error) {
				return nil
			},
		},
		{
			name:    "dial error",
			enabled: true,

			dialTimeout: time.Second,

			startServer: func(_ string) (net.Listener, error) {
				return nil, nil
			},
			acceptConnection: func(_ net.Listener) error {
				return nil
			},
			stopServer: func(_ net.Listener) error {
				return nil
			},

			wantErr: true,
		},
		{
			name:    "copy error",
			enabled: true,

			dialTimeout: time.Minute,

			startServer: func(address string) (net.Listener, error) {
				ln, err := net.Listen("tcp", address)
				if err != nil {
					return ln, err
				}

				return ln, nil
			},
			acceptConnection: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				conn, err := ln.Accept()
				if err != nil {
					if err == io.EOF {
						return nil
					}

					return err
				}

				defer conn.Close()

				_, err = conn.Write([]byte(`1111111111`))
				if err != nil {
					return err
				}

				return nil
			},
			stopServer: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				return ln.Close()
			},

			contentLength: 20,
			timeout:       time.Second,
			callback: func(r io.Reader) (err error) {
				return nil
			},

			wantErr: true,
		},
		{
			name:    "callback error",
			enabled: true,

			dialTimeout: time.Minute,

			startServer: func(address string) (net.Listener, error) {
				ln, err := net.Listen("tcp", address)
				if err != nil {
					return ln, err
				}

				return ln, nil
			},
			acceptConnection: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				conn, err := ln.Accept()
				if err != nil {
					if err == io.EOF {
						return nil
					}

					return err
				}

				defer conn.Close()

				_, err = conn.Write([]byte(`1111111111`))
				if err != nil {
					return err
				}

				return nil
			},
			stopServer: func(ln net.Listener) error {
				if ln == nil {
					return nil
				}

				return ln.Close()
			},

			contentLength: 10,
			timeout:       time.Second,
			callback: func(r io.Reader) (err error) {
				return errors.New("some error")
			},

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			addr := "127.0.0.1:" + strconv.Itoa(1024+rand.Intn(64511))

			ln, err := test.startServer(addr)
			if err != nil {
				t.Fatal(err)
			}

			errCh := make(chan error, 1)
			defer close(errCh)

			go func(ln net.Listener, accept func(net.Listener) error, ch chan<- error) {
				ch <- accept(ln)
			}(ln, test.acceptConnection, errCh)

			err = NewDownloader(addr, test.dialTimeout).Download(test.contentLength, test.timeout, test.callback)

			if stopErr := test.stopServer(ln); stopErr != nil {
				t.Fatal(stopErr)
			}

			if chErr := <-errCh; chErr != nil {
				t.Fatal(chErr)
			}

			if (err != nil) != test.wantErr {
				t.Error(err)
			}
		})
	}
}
