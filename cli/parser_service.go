package cli

import (
	"encoding/json"
	"io"
)

type ParseService struct {
	dc DownloaderCreator
	pc ParserCreator
}

func NewParseService(dc DownloaderCreator, pc ParserCreator) *ParseService {
	return &ParseService{
		dc: dc,
		pc: pc,
	}
}

func (svc *ParseService) Parse(w io.Writer, r io.Reader) (err error) {
	var (
		in = &Input{
			DialTimeout: DefaultDialTimeout,
			ReadTimeout: DefaultReadTimeout,
		}
		out = &Output{Success: true}
	)

	defer func(w io.Writer, err *error) {
		if *err == nil {
			return
		}

		*err = json.NewEncoder(w).Encode(&Output{ErrorMessage: (*err).Error()})
	}(w, &err)

	if err = json.NewDecoder(r).Decode(in); err != nil {
		return err
	}

	if err = in.Validate(); err != nil {
		return err
	}

	var (
		downloader = svc.dc(in.Address, in.DialTimeout.Duration())

		callback = func(r io.Reader) (err error) {
			if out.Nodes, err = svc.pc(in.XPathExpression).Parse(r); err != nil {
				return err
			}

			return nil
		}
	)

	if err = downloader.Download(in.ContentLength, in.ReadTimeout.Duration(), callback); err != nil {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	return enc.Encode(out)
}
