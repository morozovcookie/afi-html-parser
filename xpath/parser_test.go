package xpath

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		r          io.Reader
		expression string

		wantErr bool

		expectedNodes []string
	}{
		{
			name:    "pass",
			enabled: true,

			r: bytes.NewBufferString(`<ul>
				<li>Tag attributes</li>
				<li>Make plain text</li>
			</ul>`),
			expression: `//ul/li`,

			expectedNodes: []string{
				`<li>Tag attributes</li>`,
				`<li>Make plain text</li>`,
			},
		},
		{
			name:    "query error",
			enabled: true,

			r: bytes.NewBufferString(`<ul>
				<li>Tag attributes</li>
				<li>Make plain text</li>
			</ul>`),
			expression: string([]byte(nil)),

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			actualNodes, err := NewParser(test.expression).Parse(test.r)
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expectedNodes, actualNodes)
		})
	}
}
