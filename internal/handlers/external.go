package handlers

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
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

func (eh ExternalHandler) Run(context string) (string, error) {
	if _, ok := command.IsInPath(eh.op); !ok {
		return "", fmt.Errorf("%s: %w", eh.op, ErrCommandNotFound)
	}

	fmt.Println(parser.ArgsParse(context))
	cmd := exec.Command(eh.op, parser.ArgsParse(context)...)
	out, _ := cmd.CombinedOutput()
	//if err != nil {
	//	return "", fmt.Errorf("%s %w: %v, %s", eh.op, ErrExecutionWentWrong, err, out)
	//}

	return string(out), nil
}
