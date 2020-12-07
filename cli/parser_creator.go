package cli

import (
	ahp "github.com/morozovcookie/afihtmlparser"
)

type ParserCreator func(expression string) ahp.Parser
