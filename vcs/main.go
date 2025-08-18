package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adammwaniki/andika/vcs/handler"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to VCS Notes Manager")
	handler.PrintHelp()

	for {
		fmt.Print("vcs> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Fields(input)
		cmd := args[0]
		cmdArgs := args[1:]

		handler.HandleCommand(cmd, cmdArgs)
	}
}