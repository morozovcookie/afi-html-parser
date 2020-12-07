package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	ahp "github.com/morozovcookie/afihtmlparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseService_Parse(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		downloader func() ahp.Downloader

		parser       *ahp.MockParser
		parserInput  []interface{}
		parserOutput []interface{}

		input io.Reader

		expected func(*testing.T) string

		wantErr bool
	}{
		{
			name:    "pass",
			enabled: true,

			downloader: func() ahp.Downloader {
				return ahp.NewMockDownloaderWithParser(bytes.NewBufferString(`<li>blabla</li>`))
			},

			parser: &ahp.MockParser{},
			parserInput: []interface{}{
				bytes.NewBufferString(`<li>blabla</li>`),
			},
			parserOutput: []interface{}{
				[]string{
					`<li>blabla</li>`,
				},
				(error)(nil),
			},

			input: bytes.NewBufferString(
				`{"content-length":10,"address":"127.0.0.1:8080","xpath-expression":"//ul/li"}`),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success: true,
						Nodes: []string{
							"<li>blabla</li>",
						},
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
		{
			name:    "empty nodes list",
			enabled: true,

			downloader: func() ahp.Downloader {
				return ahp.NewMockDownloaderWithParser(bytes.NewBufferString(``))
			},

			parser: &ahp.MockParser{},
			parserInput: []interface{}{
				bytes.NewBufferString(``),
			},
			parserOutput: []interface{}{
				([]string)(nil),
				(error)(nil),
			},

			input: bytes.NewBufferString(
				`{"content-length":10,"address":"127.0.0.1:8080","xpath-expression":"//ul/li"}`),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success: true,
						Nodes:   nil,
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
		{
			name:    "decode error",
			enabled: true,

			downloader: func() ahp.Downloader {
				return nil
			},

			parser:       &ahp.MockParser{},
			parserInput:  []interface{}{},
			parserOutput: []interface{}{},

			input: bytes.NewBuffer(nil),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success:      false,
						ErrorMessage: "EOF",
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
		{
			name:    "validate error",
			enabled: true,

			downloader: func() ahp.Downloader {
				return nil
			},

			parser:       &ahp.MockParser{},
			parserInput:  []interface{}{},
			parserOutput: []interface{}{},

			input: bytes.NewBufferString(`{}`),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success:      false,
						ErrorMessage: "input validation error: zero content-length value",
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
		{
			name:    "download error",
			enabled: true,

			downloader: func() ahp.Downloader {
				var (
					downloader = &ahp.MockDownloader{}

					input = []interface{}{
						int64(10),
						time.Second,
						mock.AnythingOfType("afihtmlparser.DownloadCallback"),
					}

					output = []interface{}{
						errors.New("download error"),
					}
				)

				downloader.
					On("Download", input...).
					Return(output...)

				return downloader
			},

			parser:       &ahp.MockParser{},
			parserInput:  []interface{}{},
			parserOutput: []interface{}{},

			input: bytes.NewBufferString(
				`{"content-length":10,"address":"127.0.0.1:8080","xpath-expression":"//ul/li"}`),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success:      false,
						ErrorMessage: "download error",
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
		{
			name:    "parse error",
			enabled: true,

			downloader: func() ahp.Downloader {
				return ahp.NewMockDownloaderWithParser(nil)
			},

			parser: &ahp.MockParser{},
			parserInput: []interface{}{
				(io.Reader)(nil),
			},
			parserOutput: []interface{}{
				([]string)(nil),
				errors.New("parse error"),
			},

			input: bytes.NewBufferString(
				`{"content-length":10,"address":"127.0.0.1:8080","xpath-expression":"//ul/li"}`),

			expected: func(t *testing.T) string {
				var (
					buf = &bytes.Buffer{}

					out = &Output{
						Success:      false,
						ErrorMessage: "parse error",
					}
				)

				enc := json.NewEncoder(buf)
				enc.SetEscapeHTML(false)

				if err := enc.Encode(out); err != nil {
					t.Fatal(err)
				}

				return buf.String()
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			test.parser.
				On("Parse", test.parserInput...).
				Return(test.parserOutput...)

			var (
				downloaderCreator = func(_ string, _ time.Duration) ahp.Downloader {
					return test.downloader()
				}

				parserCreator = func(_ string) ahp.Parser {
					return test.parser
				}

				actual = &bytes.Buffer{}
			)

			err := NewParseService(downloaderCreator, parserCreator).Parse(actual, test.input)
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.Fail()
			}

			assert.Equal(t, test.expected(t), actual.String())
		})
	}
}
