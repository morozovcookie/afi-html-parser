package main

import (
	"fmt"
	"os"
	"time"

	ahp "github.com/morozovcookie/afihtmlparser"
	"github.com/morozovcookie/afihtmlparser/cli"
	"github.com/morozovcookie/afihtmlparser/tcp"
	"github.com/morozovcookie/afihtmlparser/xpath"
)

func main() {
	var (
		downloaderCreator = func(address string, timeout time.Duration) ahp.Downloader {
			return tcp.NewDownloader(address, timeout)
		}

		parserCreator = func(expression string) ahp.Parser {
			return xpath.NewParser(expression)
		}
	)

	if err := cli.NewParseService(downloaderCreator, parserCreator).Parse(os.Stdout, os.Stdin); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "parse error: %v \n", err)
	}
}
