package communicator

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
	Cd   = "cd"
)

func StartCommunication() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		trimmed := strings.SplitN(command[:len(command)-1], " ", 2)
		op, argv := trimmed[0], trimmed[1:]
		fmt.Println(validateInputQuotes(string(argv[0])))

		switch op {
		case Exit:
			os.Exit(0)
		case Echo:
			// fmt.Println(strings.Join(argv, " "))
		case Type:
			processType(argv...)
		case Pwd:
			processPwd()
		case Cd:
			processCd(argv...)
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

func validateInputQuotes(input string) string {
	var total string
	quotesIds := getInputQuotes(input)

	if len(quotesIds) == 0 {
		splited := strings.Split(input, " ")
		var content []string
		for _, v := range splited {
			if v != "" {

				content = append(content, v)
			}
		}
		return strings.Join(content, " ")
	}

	if len(quotesIds) == 2 {
		if quotesIds[0]+1 == quotesIds[1] {
			return input[:quotesIds[0]] + input[quotesIds[1]+1:]
		}
	}

	for i := 0; i < len(quotesIds); i++ {
		if i == 0 {
			total += input[:quotesIds[i]]
			continue
		}

		total += input[quotesIds[i-1]+1 : quotesIds[i]]
	}

	return total
}

func getInputQuotes(input string) []int {
	var quoted bool
	var quotes []int

	for idx, sym := range input {
		if string(sym) != "'" {
			continue
		}

		if quoted == true {
			quotes = append(quotes, idx)
			quoted = false
		} else {
			quotes = append(quotes, idx)
			quoted = true
		}
	}

	return quotes
}

func stopReadingAllowed(input string) bool {
	var quoted bool
	for _, sym := range input {
		if string(sym) != "'" {
			continue
		}

		if quoted == true {
			return true
		}

		quoted = true
	}

	return quoted == false
}

func processType(args ...string) {
	if isBuiltinOp(args[0]) {
		fmt.Printf("%s is a shell builtin\n", args[0])
	} else if path, ok := isPathCommand(args[0]); ok {
		fmt.Printf("%s is %s\n", args[0], path)
	} else {
		fmt.Println(args[0] + ": not found")
	}
}

func processCd(args ...string) {
	if len(args) == 0 {
		return
	}

	ap := args[0]

	if len(ap) == 0 {
		fmt.Println("return due to 0 length")
		return
	}

	if ap == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			ap = home
		}
	}

	fi, err := os.Stat(ap)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", ap)
		return
	}

	if !fi.IsDir() {
		return
	}

	if err := os.Chdir(ap); err != nil {
		return
	}
}

func processPwd() {
	p, _ := os.Getwd()
	fmt.Println(p)
}

func isBuiltinOp(op string) bool {
	switch op {
	case Type, Exit, Echo, Pwd, Cd:
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
