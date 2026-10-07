package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Nepsteren/CLI_FinanceTracker.git/output"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ expence-tracker > ")

		scanner.Scan()
		cmd := scanner.Text()
		switch cmd {
		case "add":
		case "update":
		case "delete":
		case "list":
		case "summary":
			//month
		case "set":
		case "export":
		case "help":
			output.PrintHelp()
		case "q":
			return
		}
	}
}
