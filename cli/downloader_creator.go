package cli

import (
	"time"

	ahp "github.com/morozovcookie/afihtmlparser"
)

type DownloaderCreator func(address string, timeout time.Duration) (downloader ahp.Downloader)
