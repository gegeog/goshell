package command

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	CmdExit = "exit"
	CmdEcho = "echo"
	CmdType = "type"
	CmdPwd  = "pwd"
	CmdCd   = "cd"
)

func IsBuiltin(op string) bool {
	switch op {
	case CmdType, CmdExit, CmdEcho, CmdPwd, CmdCd:
		return true
	}

	return false
}

func IsInPath(op string) (string, bool) {
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
