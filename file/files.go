package file

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Nepsteren/CLI_FinanceTracker.git/model"
)

var path string = "expense.json"

func GetExpenses(exp *[]model.Expence) {
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

func AddExpence(expences *[]model.Expence, desc string, amount float64) {
	exp := append(*expences, model.Expence{
		Id:     NextId(expences),
		Desc:   desc,
		Amount: amount,
		Date:   time.Now(),
	})

	WriteFile(&exp)
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
