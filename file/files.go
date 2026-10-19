package file

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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

func Load(exp *[]model.Expence) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)

	err = json.Unmarshal(data, &exp)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadBudget(budget *float64) {
	f, err := os.OpenFile(pathBudget, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)

	err = json.Unmarshal(data, &budget)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFileB(budget *float64) {
	data, err := json.MarshalIndent(&budget, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(pathBudget, data, 0644)
}

func WriteFile(exp *[]model.Expence) {
	data, err := json.MarshalIndent(exp, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(path, data, 0644)
}

func NextId(exp *[]model.Expence) int {
	id := len(*exp)
	return id
}

func SearchExp(exp *[]model.Expence, id int) int {
	for i := 0; i < len(*exp); i++ {
		if (*exp)[i].Id == id {
			return i
		}
	}
	return -1
}

func Export(epx *[]model.Expence) {
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
			fmt.Printf("incorrect budget - %d\n", budget)
			return -1
		}
		budget = bud
		WriteFileB(&budget)
	}
	return budget
}

func AddExp(expences *[]model.Expence, desc string, amount float64) {
	id := NextId(expences)
	exp := append(*expences, model.Expence{
		Id:     id,
		Desc:   desc,
		Amount: amount,
		Date:   time.Now(),
	})

	fmt.Println("expense added successfully -", id)
	WriteFile(&exp)
}

func UpdateChecksFlags(flags []string) (int, string, float64) {
	var desk string
	var amount float64
	var Id int

	var foundDesk bool
	var foundAmount bool
	var foundId bool

	for i := 1; i < len(flags); i++ {
		if flags[i] == "--description" && len(flags) > i+1 {
			desk = flags[i+1]
			foundDesk = true
		} else if flags[i] == "--amount" && len(flags) > i+1 {
			amount, _ = strconv.ParseFloat(flags[i+1], 64)
			foundAmount = true
		} else if flags[i] == "--id" && len(flags) > i+1 {
			foundId = true
			id, err := strconv.Atoi(flags[i+1])
			if err != nil {
				log.Fatal(err)
			}
			Id = id
		}

	}

	if foundId {
		if foundAmount && foundDesk {
			return Id, desk, amount
		}

		if foundAmount {
			return Id, "", amount
		}

		if foundDesk {
			return Id, desk, -1
		}
	}

	return -1, "", 0
}

func UpdateExp(exp *[]model.Expence, id int, desk string, amount float64) {
	for i := 0; i < len((*exp)); i++ {
		if (*exp)[i].Id == id {
			if desk != "" {
				(*exp)[i].Desc = desk
			}
			if (*exp)[i].Amount == -1 {
				(*exp)[i].Amount = amount
			}
		}
	}
	WriteFile(exp)

}

func DeleteExp(exp *[]model.Expence, id int) {
	SearchExp(exp, id)
	if id > -1 {
		*exp = append((*exp)[0:id], (*exp)[id+1:]...)
		fmt.Println("expense deleted successfully -", id)
		WriteFile(exp)
	} else {
		fmt.Println("id not exist")
	}
}

func ListExp(exp *[]model.Expence) {
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

func ListSummary(exp *[]model.Expence, month string, withMonth bool, busget float64) {
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

func DeleteChecksFlags(flags []string) int {
	if flags[1] == "--id" && len(flags) > 2 {
		id, err := strconv.Atoi(flags[2])
		if err != nil {
			fmt.Println("incorrect id")
			return -1
		}
		return id
	}
	return -1
}

func AddCheckFlags(flags []string) (string, float64) {
	var desk string
	var amount float64

	var foundDesk bool
	var foundAmount bool

	for i := 1; i < len(flags); i++ {
		if flags[i] == "--description" && len(flags) > i+1 {
			desk = flags[i+1]
			foundDesk = true
		} else if flags[i] == "--amount" && len(flags) > i+1 {
			amount, _ = strconv.ParseFloat(flags[i+1], 64)
			foundAmount = true
		}

	}

	if foundAmount && foundDesk {
		return desk, amount
	}

	if foundAmount {
		return "", amount
	}

	if foundDesk {
		return desk, 0
	}

	return "", 0
}
