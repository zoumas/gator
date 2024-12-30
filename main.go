package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	s, err := newState()
	if err != nil {
		log.Fatalln(err)
	}
	_ = s

	commandManager := newCommands()

	if len(os.Args) < 2 {
		log.Fatalln("expected command")
	}
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	cmd := command{
		name: os.Args[1],
		args: args,
	}
	err = commandManager.run(s, cmd)
	if err != nil {
		log.Fatalln(err)
	}
}
