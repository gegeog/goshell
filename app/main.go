package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/handlers"
	"github.com/codecrafters-io/shell-starter-go/internal/router"
	"github.com/codecrafters-io/shell-starter-go/internal/shell"
)

func main() {
	r := router.New()
	r.Handle(command.CmdExit, handlers.ExitHandler{})
	r.Handle(command.CmdType, handlers.TypeHandler{})
	r.Handle(command.CmdPwd, handlers.PwdHandler{})
	r.Handle(command.CmdCd, handlers.CdHandler{})
	r.Handle(command.CmdEcho, handlers.EchoHandler{})

	if err := shell.ListenAndServe(r); err != nil {
		// log.Fatal(err)
		os.Exit(1)
	}
}
