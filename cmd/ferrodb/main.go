package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ferrodb/internal/engine"
)

func main() {
	fmt.Println("FerroDB v0.1")
	fmt.Println("Type HELP for commands")

	eng := engine.New()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.ToUpper(line) == "EXIT" {
			fmt.Println("Bye 👋")
			break
		}

		result := eng.Execute(line)
		fmt.Println(result)
	}
}
