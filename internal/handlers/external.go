package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

var ErrCommandNotFound = errors.New("command not found")

type ExternalHandler struct {
	op string
}

func NewExternal(op string) ExternalHandler {
	return ExternalHandler{
		op,
	}
}

func (eh ExternalHandler) Run(context []string) (string, error) {
	if _, ok := command.IsInPath(eh.op); !ok {
		return "", fmt.Errorf("%s: %w", eh.op, ErrCommandNotFound)
	}

	var buf bytes.Buffer

	cmd := exec.Command(eh.op, context...)
	cmd.Stderr = &buf
	out, err := cmd.Output()

	if err != nil && buf.Len() > 0 {
		return string(out), fmt.Errorf(buf.String())
	}

	return string(out), nil
}
