package handlers

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

var ErrCommandNotFound = errors.New("command not found")
var ErrExecutionWentWrong = errors.New("execution went wrong")

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

	cmd := exec.Command(eh.op, context...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), fmt.Errorf("%s %w: %v, %s", eh.op, ErrExecutionWentWrong, err, out)
	}

	return string(out), nil
}
