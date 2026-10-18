package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Nepsteren/CLI_FinanceTracker.git/file"
	"github.com/Nepsteren/CLI_FinanceTracker.git/model"
	"github.com/Nepsteren/CLI_FinanceTracker.git/output"
)

var exp []model.Expence

func main() {
	file.Load(&exp)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		output.Start()

		scanner.Scan()
		input := scanner.Text()

		flags := strings.Split(input, " ")
		cmd := flags[0]

		switch cmd {
		case "add":
			desk, amount := file.AddCheckFlags(flags)
			file.AddExp(&exp, desk, amount)
		case "update":
			id, desc, amount := file.UpdateChecksFlags(flags)
			if id != -1 {
				file.UpdateExp(&exp, id, desc, amount)
			} else {
				fmt.Printf("incorrect id - %d\n", id)
			}

		case "delete":
			id := file.DeleteChecksFlags(flags)
			if id == -1 {
				fmt.Println("not such a task")
				break
			}
			file.DeleteExp(&exp, id)
		case "list":
			file.ListExp(&exp)
		case "summary":
			month, withMonth := file.SummaryChecksFlags(flags)
			file.ListSummary(&exp, month, withMonth)
		case "set":
		case "export":
		case "help":
			output.PrintHelp()
		case "q":
			return
		default:
			output.InvalidCmd(cmd)
		}

	}
}
