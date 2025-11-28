package handlers

import (
	"errors"
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

var ErrNotFound = errors.New("not found")

type TypeHandler struct{}

func (th TypeHandler) Run(input string) (string, error) {
	if command.IsBuiltin(input) {
		return fmt.Sprintf("%s is a shell builtin", input), nil
	}

	if path, ok := command.IsInPath(input); ok {
		return fmt.Sprintf("%s is %s", input, path), nil
	}

	return "", fmt.Errorf("%s: %w", input, ErrNotFound)
}
