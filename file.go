package afihtmlparser

import (
	"io"
	"time"
)

type DownloadCallback func(r io.Reader) (err error)

type Downloader interface {
	Download(contentLength int64, timeout time.Duration, callback DownloadCallback) (err error)
}

type Parser interface {
	Parse(r io.Reader) ([]string, error)
}
