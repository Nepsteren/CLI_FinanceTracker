package file

import (
	"encoding/json"
	"fmt"
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

func SearchExp(exp *[]model.Expence, id int) int {
	for i := 0; i < len(*exp); i++ {
		if (*exp)[i].Id == id {
			return i
		}
	}
	return -1
}

func AddExp(expences *[]model.Expence, desc string, amount float64) {
	exp := append(*expences, model.Expence{
		Id:     NextId(expences),
		Desc:   desc,
		Amount: amount,
		Date:   time.Now(),
	})

	WriteFile(&exp)
}

func DeleteExp(exp *[]model.Expence, id int) {
	SearchExp(exp, id)
	if id > -1 {
		*exp = append((*exp)[0:id], (*exp)[id+1:]...)
		WriteFile(exp)
	} else {
		fmt.Println("id not exist")
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
