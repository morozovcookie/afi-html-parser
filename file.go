package afihtmlparser

import (
	"github.com/stretchr/testify/mock"
	"io"
	"time"
)

type DownloadCallback func(r io.Reader) (err error)

type Downloader interface {
	Download(contentLength int64, timeout time.Duration, callback DownloadCallback) (err error)
}

type MockDownloader struct {
	mock.Mock
}

func (d *MockDownloader) Download(contentLength int64, timeout time.Duration, callback DownloadCallback) (err error) {
	return d.Called(contentLength, timeout, callback).Error(0)
}

type MockDownloaderWithParser struct {
	r io.Reader
}

func NewMockDownloaderWithParser(r io.Reader) *MockDownloaderWithParser {
	return &MockDownloaderWithParser{
		r: r,
	}
}

func (d *MockDownloaderWithParser) Download(_ int64, _ time.Duration, callback DownloadCallback) (err error) {
	return callback(d.r)
}

type Parser interface {
	Parse(r io.Reader) (nodes []string, err error)
}

type MockParser struct {
	mock.Mock
}

func (p *MockParser) Parse(r io.Reader) (nodes []string, err error) {
	args := p.Called(r)

	return args.Get(0).([]string), args.Error(1)
}
