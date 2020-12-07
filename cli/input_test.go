package cli

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		d Duration

		wantErr  bool
		expected []byte
	}{
		{
			name:    "pass",
			enabled: true,

			d: Duration(time.Second),

			expected: []byte(`"1s"`),
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			actual, err := test.d.MarshalJSON()
			if (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		in []byte

		expected Duration

		wantErr bool
	}{
		{
			name:    "unmarshal from number",
			enabled: true,

			in: []byte(strconv.FormatInt(int64(time.Second), 10)),

			expected: Duration(time.Second),
		},
		{
			name:    "unmarshal from string",
			enabled: true,

			in: []byte(`"1s"`),

			expected: Duration(time.Second),
		},
		{
			name:    "unmarshal error",
			enabled: true,

			wantErr: true,
		},
		{
			name:    "parse duration error",
			enabled: true,

			in: []byte(`"1"`),

			wantErr: true,
		},
		{
			name:    "invalid duration",
			enabled: true,

			in: []byte{0x7B, 0x7D}, // {}

			wantErr: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			var actual Duration
			if err := actual.UnmarshalJSON(test.in); (err != nil) != test.wantErr {
				t.Error(err)
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDuration_Duration(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		duration Duration

		expected time.Duration
	}{
		{
			name:    "pass",
			enabled: true,

			duration: Duration(time.Second),

			expected: time.Second,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			assert.Equal(t, test.expected, test.duration.Duration())
		})
	}
}

func TestInput_Validate(t *testing.T) {
	tt := []struct {
		name    string
		enabled bool

		input *Input

		wantErr  bool
		expected error
	}{
		{
			name:    "pass with host:port address",
			enabled: true,

			input: &Input{
				ContentLength:   10,
				Address:         "127.0.0.1:8080",
				XPathExpression: "//ul/li",
			},
		},
		{
			name:    "pass with hostname:port address",
			enabled: true,

			input: &Input{
				ContentLength:   10,
				Address:         "mydomain.zone",
				XPathExpression: "//ul/li",
			},
		},
		{
			name:    "zero content length",
			enabled: true,

			input: &Input{},

			wantErr:  true,
			expected: ErrZeroContentLengthValue,
		},
		{
			name:    "empty address",
			enabled: true,

			input: &Input{
				ContentLength: 10,
			},

			wantErr:  true,
			expected: ErrEmptyAddress,
		},
		{
			name:    "invalid host:port address",
			enabled: true,

			input: &Input{
				ContentLength: 10,
				Address:       "256.789.320.752:8135135368",
			},

			wantErr:  true,
			expected: ErrInvalidAddress,
		},
		{
			name:    "invalid hostname:port address",
			enabled: true,

			input: &Input{
				ContentLength: 10,
				Address:       "gsfdsfdfd%@#fdfaf",
			},

			wantErr:  true,
			expected: ErrInvalidAddress,
		},
		{
			name:    "empty xpath expression",
			enabled: true,

			input: &Input{
				ContentLength: 10,
				Address:       "127.0.0.1:8080",
			},

			wantErr:  true,
			expected: ErrEmptyXPathExpression,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.SkipNow()
			}

			actual := test.input.Validate()
			if (actual != nil) != test.wantErr {
				t.Error(actual)
				t.FailNow()
			}

			if test.wantErr {
				assert.NotNil(t, actual)
				assert.EqualError(t, actual, test.expected.Error())
				return
			}

			assert.Nil(t, actual)
		})
	}
}
