package handlers

import (
	"errors"
	"fmt"
	"os"
)

var ErrUnableOpenHomeDirectory = errors.New("unable to open home directory")
var ErrUnableToChangeDirectory = errors.New("unable to go to provided dir")
var ErrPathIsNotDirectory = errors.New("provided path is not directory")

var ErrNoSuchFileOrDirectory = errors.New("No such file or directory")

type CdHandler struct{}

func (ch CdHandler) Run(args []string) (string, error) {
	var path string

	if len(args) == 0 || args[0] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", ErrUnableOpenHomeDirectory
		}

		path = home
	} else {
		path = args[0]
	}

	fi, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cd: %s: %w", path, ErrNoSuchFileOrDirectory)
	}

	if !fi.IsDir() {
		return "", ErrPathIsNotDirectory
	}

	if err := os.Chdir(path); err != nil {
		return "", ErrUnableToChangeDirectory
	}

	return "", nil
}
