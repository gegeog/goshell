package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	Exit = "exit"
	Echo = "echo"
	Type = "type"
	Pwd  = "pwd"
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
			} else if path, ok := isPathCommand(argv[0]); ok {
				fmt.Printf("%s is %s\n", argv[0], path)
			} else {
				fmt.Println(argv[0] + ": not found")
			}
		case Pwd:
			p, _ := os.Getwd()
			if err == nil {
				fmt.Println(p)
			}
		default:
			if _, ok := isPathCommand(op); !ok {
				fmt.Println(op + ": command not found")
			} else {
				cmd := exec.Command(op, argv...)
				out, err := cmd.CombinedOutput()
				if err == nil {
					fmt.Print(string(out))
				}
			}
		}
	}
}

func isBuiltinOp(op string) bool {
	switch op {
	case Type, Exit, Echo, Pwd:
		return true
	}

	return false
}

func isExecutable(info fs.FileInfo) bool {
	// Windows logic: Check file extension
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(info.Name()))
		switch ext {
		case ".exe", ".bat", ".cmd", ".ps1":
			return true
		default:
			return false
		}
	}

	// Unix logic: Check execute bit (0111 octal = 73 decimal)
	// We check if at least one execute bit is set (Owner, Group, or Other)
	return info.Mode()&0111 != 0
}

func isPathCommand(op string) (string, bool) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))

	for _, dir := range dirs {

		fullPath := filepath.Join(dir, op)
		if fileInfo, err := os.Stat(fullPath); err == nil {
			if !fileInfo.IsDir() && isExecutable(fileInfo) {
				return fullPath, true
			} else {
				continue
			}
		}
	}

	return "", false
}
