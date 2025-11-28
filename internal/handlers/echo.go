package handlers

import "github.com/codecrafters-io/shell-starter-go/internal/parser"

type EchoHandler struct{}

func (eh EchoHandler) Run(s string) (string, error) {
	return parser.EchoParse(s), nil
}
