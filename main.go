package main

import (
	"dara/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, welcome to the dara repl!\n",
		user.Username)
	fmt.Printf("Type in commands below to evaluate your code.\n")
	repl.Start(os.Stdin, os.Stdout)
}
