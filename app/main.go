package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Exit = "exit"
	Echo = "echo"
	Type = "type"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		trimmed := strings.Split(command[:len(command)-1], " ")
		op, argv := trimmed[0], trimmed[1:]

		switch op {
		case Exit:
			os.Exit(0)
		case Echo:
			fmt.Println(strings.Join(argv, " "))
		case Type:
			if isBuiltinOp(argv[0]) {
				fmt.Printf("%s is a shell builtin\n", argv[0])
			} else {
				fmt.Println(argv[0] + ": not found")
			}
		default:
			fmt.Println(op + ": command not found")
		}
	}
}

func isBuiltinOp(op string) bool {
	switch op {
	case Type, Exit, Echo:
		return true
	}

	return false
}
