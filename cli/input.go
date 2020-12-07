package cli

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

const (
	DefaultDialTimeout = Duration(time.Second)
	DefaultReadTimeout = Duration(time.Second)
)

var ErrInvalidDuration = errors.New("invalid duration")

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
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
	DialTimeout     Duration `json:"dial-timeout"`
	ReadTimeout     Duration `json:"read-timeout"`
}

var (
	ErrZeroContentLengthValue = errors.New("input validation error: zero content-length value")
	ErrEmptyAddress         = errors.New("input validation error: empty address")
	ErrInvalidAddress         = errors.New("input validation error: invalid address")
	ErrEmptyXPathExpression   = errors.New("input validation error: empty xpath expression")
)

func (i Input) Validate() (err error) {
	if i.ContentLength <= 0 {
		return ErrZeroContentLengthValue
	}

	if err = validateAddress(i.Address); err != nil {
		return err
	}

	if i.XPathExpression == "" {
		return ErrEmptyXPathExpression
	}

	return nil
}

const (
	HostPortRegex = `(?m)^((((25[0-5])|(2[0-4]\d{1})|([0-1]?\d{1,2}))\.){3}((25[0-5])|(2[0-4]\d{1})|` +
		`([0-1]?\d{1,2})){1}(:((6553[0-5])|(655[0-2]\d{1})|(65[0-4]\d{2})|(6[0-4]\d{3})|([1-5]\d{4})|` +
		`([1-9]\d{3})|([1-9]\d{2})|([1-9]\d{1})|([1-9])))?)$`
	HostnamePortRegex = `(?m)^(((([\d\w]|[\d\w][\d\w\-]*[\d\w])\.)*([\d\w]|[\d\w][\d\w\-]*[\d\w]))` +
		`(:((6553[0-5])|(655[0-2]\d{1})|(65[0-4]\d{2})|(6[0-4]\d{3})|([1-5]\d{4})|([1-9]\d{3})|` +
		`([1-9]\d{2})|([1-9]\d{1})|([1-9])))?)$`
)

func validateAddress(s string) (err error) {
	if s == "" {
		return ErrEmptyAddress
	}

	if ok := regexp.MustCompile(HostPortRegex).MatchString(s); ok {
		return nil
	}

	if ok := regexp.MustCompile(HostnamePortRegex).MatchString(s); ok {
		return nil
	}

	return ErrInvalidAddress
}
