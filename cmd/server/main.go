package main

import (
	"fmt"
	"github.com/kenlabs/pandofg/cmd/server/command"
)

func main() {
	rootCmd := command.NewRoot()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Exit with error.")
	}
}
