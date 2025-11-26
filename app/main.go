package main

import (
	"github.com/codecrafters-io/shell-starter-go/internal/communicator"
)

func main() {
	if err := communicator.StartCommunication(); err != nil {
		panic(err)
	}
}
