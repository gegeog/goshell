package handlers

import (
	"errors"
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

var ErrNotFound = errors.New("not found")
var ErrNilArgs = errors.New("wrong arguments")

type TypeHandler struct{}

func (th TypeHandler) Run(args []string) (string, error) {
	if len(args) < 1 {
		return "", ErrNilArgs
	}

	if command.IsBuiltin(args[0]) {
		return fmt.Sprintf("%s is a shell builtin", args[0]), nil
	}

	if path, ok := command.IsInPath(args[0]); ok {
		return fmt.Sprintf("%s is %s", args[0], path), nil
	}

	return "", fmt.Errorf("%s: %w", args[0], ErrNotFound)
}
