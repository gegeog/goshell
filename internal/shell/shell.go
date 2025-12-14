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

		op, parseinfo := parser.Parse(command)

		if op == "" {
			continue
		}

		out, err := r.Run(op, parseinfo.Arguments)

		if errors.Is(err, handlers.ErrShellExit) {
			return err
		}

		if errors.Is(err, handlers.ErrNoSuchFileOrDirectory) {
			writeError(err, parseinfo)
			continue
		}

		if errors.Is(err, handlers.ErrCommandNotFound) {
			writeError(err, parseinfo)
			continue
		}

		if errors.Is(err, handlers.ErrNotFound) {
			writeError(err, parseinfo)
			continue
		}

		/*
			if no error but error output file provided
			should create empty file
		*/
		writeLine(out, parseinfo)
		writeError(err, parseinfo)
	}
}

func writeLine(s string, info parser.ParsedInfo) {
	if s == "" && info.OutputRedirectsRest == nil && info.OutputRedirectsNew == nil {
		return
	}

	if info.OutputRedirectsRest == nil && info.OutputRedirectsNew == nil {
		fmt.Fprintln(os.Stdout, strings.TrimRight(s, "\n"))
		return
	}

	for _, outPath := range info.OutputRedirectsNew {
		file, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0644)
		_, _ = file.WriteString(s)
		_ = file.Close()
	}

	for _, outPath := range info.OutputRedirectsRest {
		file, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		_, _ = file.WriteString(s + "\n")
		_ = file.Close()
	}
}

func writeError(err error, info parser.ParsedInfo) {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	for _, errPath := range info.ErrRedirectNew {
		file, _ := os.OpenFile(errPath, os.O_WRONLY|os.O_CREATE, 0644)
		_, _ = file.WriteString(msg)
		_ = file.Close()
	}

	for _, errPath := range info.ErrRedirectRest {
		file, _ := os.OpenFile(errPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		_, _ = file.WriteString(msg)
		_ = file.Close()
	}
}
