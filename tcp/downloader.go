package tcp

import (
	"bytes"
	"io"
	"net"
	"time"

	ahp "github.com/morozovcookie/afihtmlparser"
)

type Downloader struct {
	address string
	timeout time.Duration
}

func NewDownloader(address string, timeout time.Duration) *Downloader {
	return &Downloader{
		address: address,
		timeout: timeout,
	}
}

func (d *Downloader) Download(contentLength int64, timeout time.Duration, callbackFn ahp.DownloadCallback) (err error) {
	conn, err := net.DialTimeout("tcp", d.address, d.timeout)
	if err != nil {
		return err
	}

	defer conn.Close()

	if err = conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	buf := &bytes.Buffer{}
	if _, err = io.CopyN(buf, conn, contentLength); err != nil {
		return err
	}

	return callbackFn(buf)
}
