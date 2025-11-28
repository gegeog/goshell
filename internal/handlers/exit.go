package handlers

import "errors"

var ErrShellExit = errors.New("shell exit")

type ExitHandler struct{}

func (eh ExitHandler) Run(string) (string, error) {
	return "", ErrShellExit
}
