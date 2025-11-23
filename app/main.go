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
			os.Exit(1)
		case Echo:
			fmt.Println(strings.Join(argv, " "))
		default:
			fmt.Println(op + ": command not found")
		}
	}
}
