package handlers

import (
	"strings"
)

type EchoHandler struct{}

func (eh EchoHandler) Run(args []string) (string, error) {
	return strings.Join(args, " "), nil
}
