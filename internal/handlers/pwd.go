package handlers

import (
	"os"
)

type PwdHandler struct{}

func (ph PwdHandler) Run(string) (string, error) {
	return os.Getwd()
}
