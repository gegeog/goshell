package shell

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/handlers"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
	"github.com/codecrafters-io/shell-starter-go/internal/router"
)

func ListenAndServe(r router.Router) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading input: %v\n", err))
		}

		op, argv, outputPaths, errorPaths := parser.Parse(command)

		if op == "" {
			continue
		}

		out, err := r.Run(op, argv)

		if errors.Is(err, handlers.ErrShellExit) {
			return err
		}

		if errors.Is(err, handlers.ErrNoSuchFileOrDirectory) {
			writeLine(err.Error(), errorPaths)
			continue
		}

		if errors.Is(err, handlers.ErrCommandNotFound) {
			writeLine(err.Error(), errorPaths)
			continue
		}

		if errors.Is(err, handlers.ErrNotFound) {
			writeLine(err.Error(), errorPaths)
			continue
		}

		/*
			if no error but error output file provided
			should create empty file
		*/
		if errorPaths != nil {
			writeLine(out, outputPaths)
			if err != nil {
				writeLine(err.Error(), errorPaths)
				continue
			}

			writeLine("", errorPaths)
			continue
		}

		writeLine(out, outputPaths)
		if err != nil {
			writeLine(err.Error(), errorPaths)
		}
	}
}

func writeLine(s string, filePaths []string) {
	if s == "" && filePaths == nil {
		return
	}

	if filePaths == nil {
		fmt.Fprintln(os.Stdout, strings.TrimRight(s, "\n"))
		return
	}

	for _, v := range filePaths {
		_ = os.WriteFile(v, []byte(s), 0644)
	}
}
