package file

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nepsteren/CLI_FinanceTracker.git/model"
)

var path string = "expense.json"
var pathBudget string = "budget.json"
var pathExport string = "expences.csv"

func Load(exp *[]model.Expense) {
	loadJSON(path, exp)
}

func LoadBudget(budget *float64) {
	loadJSON(pathBudget, budget)
}

func loadJSON(filePath string, target any) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Printf("file read error: %v", err)
		return
	}

	if len(data) == 0 {
		return
	}

	if err := json.Unmarshal(data, target); err != nil {
		log.Printf("json unmarshal error: %v", err)
	}
}

func saveJSON(filePath string, data any) {
	content, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		log.Printf("file write error: %v", err)
	}
}

func WriteFileB(budget *float64) {
	saveJSON(pathBudget, budget)
}

func WriteFile(exp *[]model.Expense) {
	saveJSON(path, exp)
}

func NextId(exp *[]model.Expense) int {
	maxID := -1

	for _, e := range *exp {
		if e.Id > maxID {
			maxID = e.Id
		}
	}

	return maxID + 1
}

func SearchExp(exp *[]model.Expense, id int) int {
	for i := 0; i < len(*exp); i++ {
		if (*exp)[i].Id == id {
			return i
		}
	}
	return -1
}

func Export(epx *[]model.Expense) {
	f, err := os.Create(pathExport)

	if err != nil {
		log.Print(err)
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer func() {
		w.Flush()
		if err := w.Error(); err != nil {
			log.Printf("csv flush error: %v", err)
		}
	}()

	if err := w.Write([]string{"Id", "Description", "Amount", "Date"}); err != nil {
		log.Print(err)
		return
	}

	for _, e := range *epx {
		row := []string{
			strconv.Itoa(e.Id),
			e.Desc,
			strconv.FormatFloat(e.Amount, 'f', 2, 64),
			e.Date.Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(row); err != nil {
			log.Print(err)
			return
		}

	}
	if err := w.Error(); err != nil {
		log.Print(err)
		return
	}

}

func SetBudget(flags []string) float64 {
	budget := 0.0
	if len(flags) >= 2 {
		bud, err := strconv.ParseFloat(flags[1], 64)
		if err != nil {
			fmt.Printf("incorrect budget - %.2f\n", budget)
			return -1
		}
		budget = bud
		WriteFileB(&budget)
	}
	return budget
}

func AddExp(expences *[]model.Expense, desc string, amount float64) {
	id := NextId(expences)
	*expences = append(*expences, model.Expense{
		Id:     id,
		Desc:   desc,
		Amount: amount,
		Date:   time.Now(),
	})

	fmt.Println("expense added successfully -", id)
	WriteFile(expences)
}

func UpdateChecksFlags(args []string) (model.ExpenseUpdate, error) {
	flags, err := ParseFlags(args)

	if err != nil {
		return model.ExpenseUpdate{}, err
	}

	if err := ValidateFlags(flags, "--id", "--description", "--amount"); err != nil {
		return model.ExpenseUpdate{}, err
	}

	idValue, exists := flags["--id"]
	if !exists {
		return model.ExpenseUpdate{}, fmt.Errorf("id is required")
	}

	id, err := strconv.Atoi(idValue)
	if err != nil || id < 0 {
		return model.ExpenseUpdate{}, fmt.Errorf("incorrect id")
	}

	update := model.ExpenseUpdate{
		ID: id,
	}

	if desc, exists := flags["--description"]; exists {
		if strings.TrimSpace(desc) == "" {
			return model.ExpenseUpdate{}, fmt.Errorf("description cannot be empty")
		}

		update.Description = &desc
	}

	if amountValue, exists := flags["--amount"]; exists {
		amount, err := strconv.ParseFloat(amountValue, 64)
		if err != nil || amount < 0 {
			return model.ExpenseUpdate{}, fmt.Errorf("incorrect amount")
		}

		update.Amount = &amount
	}
	if update.Description == nil && update.Amount == nil {
		return model.ExpenseUpdate{}, fmt.Errorf("nothing to update")
	}

	return update, nil
}

func UpdateExp(expenses *[]model.Expense, update model.ExpenseUpdate) error {
	index := SearchExp(expenses, update.ID)

	if index == -1 {
		return fmt.Errorf("expense with id %d not found", update.ID)
	}
	if update.Description != nil {
		(*expenses)[index].Desc = *update.Description
	}

	if update.Amount != nil {
		(*expenses)[index].Amount = *update.Amount
	}

	WriteFile(expenses)

	return nil
}

func DeleteExp(exp *[]model.Expense, id int) {
	index := SearchExp(exp, id)

	if index == -1 {
		fmt.Println("expence not found")
		return
	}

	*exp = append((*exp)[:index], (*exp)[index+1:]...)
	fmt.Println("expense deleted successfully -", id)
	WriteFile(exp)
}

func ListExp(exp *[]model.Expense) {
	fmt.Printf("%-4s %-20s %10s %-12s\n", "ID", "DESC", "AMOUNT", "DATE")
	fmt.Println(strings.Repeat("-", 80))
	for _, e := range *exp {
		fmt.Printf("%-4d %-20s %10.2f %-12s\n", e.Id, e.Desc, e.Amount, e.Date.Format("2006-01-02 15:04"))
	}
}

func SummaryChecksFlags(flags []string) (string, bool) {
	if len(flags) > 2 {
		if flags[1] == "--month" {
			return flags[2], true
		}
	} else {
		return "", false
	}
	return "", false
}

func ListSummary(exp *[]model.Expense, month string, withMonth bool, busget float64) {
	sum := 0.0
	if !withMonth {
		for i := 0; i < len((*exp)); i++ {
			sum += (*exp)[i].Amount
		}
		if busget < sum {
			fmt.Printf("ваш бюджет на месяц был превышен!! превышение на %0.2f$\n", sum-busget)
		}
		fmt.Printf("total expenses: $%.2f\n", sum)
	} else {
		for i := 0; i < len((*exp)); i++ {
			if (*exp)[i].Date.Format("1") == month {
				sum += (*exp)[i].Amount
			}

		}

		fmt.Printf("total expenses %s: $%.2f\n", month, sum)
	}

}

func DeleteChecksFlags(args []string) (int, error) {
	flags, err := ParseFlags(args)
	if err != nil {
		return 0, err
	}

	if err := ValidateFlags(flags, "--id"); err != nil {
		return 0, err
	}

	idValue, exists := flags["--id"]
	if !exists {
		return 0, fmt.Errorf("id is required")
	}

	id, err := strconv.Atoi(idValue)
	if err != nil || id < 0 {
		return 0, fmt.Errorf("incorrect id")
	}

	return id, nil
}

func AddCheckFlags(args []string) (string, float64, error) {
	flags, err := ParseFlags(args)
	if err != nil {
		fmt.Println("incorrect arguments:", err)
		return "", 0, err
	}

	if err := ValidateFlags(flags, "--description", "--amount"); err != nil {
		return "", 0, err
	}

	desk, exists := flags["--description"]
	if !exists || strings.TrimSpace(desk) == "" {
		return "", 0, fmt.Errorf("description is required")
	}

	amountValue, exists := flags["--amount"]
	if !exists {
		return "", 0, fmt.Errorf("amount is required")
	}

	amount, err := strconv.ParseFloat(amountValue, 64)
	if err != nil || amount < 0 {
		return "", 0, fmt.Errorf("incorrect amount")
	}

	return desk, amount, nil

	// amount := 0.0
	// if value, exists := flags["--amount"]; exists {
	// 	parsed, err := strconv.ParseFloat(value, 64)
	// 	if err != nil || parsed < 0 {
	// 		fmt.Println("incorrect amount")
	// 		return "", 0, err

	// 	}
	// 	amount = parsed
	// }

	// if _, exists := flags["--amount"]; !exists {
	// 	fmt.Println("amount is required")
	// 	return "", 0, err
	// }

	// return desk, amount, nil
}
