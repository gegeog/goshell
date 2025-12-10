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

		op, argv, output, redirectMode := parser.Parse(command)

		if op == "" {
			continue
		}

		out, err := r.Run(op, argv)

		if errors.Is(err, handlers.ErrNoSuchFileOrDirectory) {
			writeLine(err.Error(), "")
			continue
		}

		if errors.Is(err, handlers.ErrCommandNotFound) {
			writeLine(err.Error(), "")
			continue
		}

		if errors.Is(err, handlers.ErrNotFound) {
			writeLine(err.Error(), "")
			continue
		}

		if errors.Is(err, handlers.ErrShellExit) {
			return err
		}

		if errors.Is(err, handlers.ErrExecutionWentWrong) {
			writeLine(out, "")
			continue
		}

		if redirectMode == parser.ErrorRedirect {
			writeLine(out, "")
			writeLine(err.Error(), output)
		} else {
			writeLine(out, "")
		}
	}
}

func writeLine(s string, output string) {
	if s == "" {
		return
	}

	if output != "" {
		_ = os.WriteFile(output, []byte(s), 0644)
	} else {
		fmt.Fprintln(os.Stdout, strings.TrimRight(s, "\n"))
	}
}
