package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Nepsteren/CLI_FinanceTracker.git/file"
	"github.com/Nepsteren/CLI_FinanceTracker.git/model"
	"github.com/Nepsteren/CLI_FinanceTracker.git/output"
)

var exp []model.Expense
var budget float64

func main() {
	file.Load(&exp)
	file.LoadBudget(&budget)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		output.Start()

		scanner.Scan()
		input := scanner.Text()

		flags, err := parseInput(input)

		if err != nil {
			fmt.Println("incorrect command:", err)
			continue
		}

		if len(flags) == 0 {
			continue
		}

		cmd := flags[0]

		switch cmd {
		case "add":
			desk, amount, err := file.AddCheckFlags(flags)
			if err != nil {
				fmt.Println("incorrect arguments:", err)
				break
			}
			file.AddExp(&exp, desk, amount)
		case "update":
			update, err := file.UpdateChecksFlags(flags)
			if err != nil {
				fmt.Println("incorrect arguments:", err)
				break
			}

			if err := file.UpdateExp(&exp, update); err != nil {
				fmt.Println("incorrect arguments:", err)
				break
			}
		case "delete":
			id, err := file.DeleteChecksFlags(flags)
			if err != nil {
				fmt.Println("incorrect arguments:", err)
				break
			}

			file.DeleteExp(&exp, id)

		case "list":
			file.ListExp(&exp)
		case "summary":
			month, withMonth := file.SummaryChecksFlags(flags)
			file.ListSummary(&exp, month, withMonth, budget)
		case "set":
			budget = file.SetBudget(flags)
		case "export":
			file.Export(&exp)
		case "help":
			output.PrintHelp()
		case "q":
			return
		default:
			output.InvalidCmd(cmd)
		}

	}
}
