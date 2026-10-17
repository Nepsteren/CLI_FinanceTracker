package output

import "fmt"

func Start() {
	fmt.Print("$ expence-tracker > ")
}

func InvalidCmd(cmd string){
	fmt.Println("unknown command -", cmd)
}

func PrintHelp() {
	fmt.Println("Список команд:")
	fmt.Println("add (description, amount) - добавить расход")
	fmt.Println("delete (id) - удалить расход")
	fmt.Println("update (id, description) - обновить расход")
	fmt.Println("list - вывести список всех расходов")
	fmt.Println("summary (month | category) - вывести сводку всех расходов | сводку за месяц | сводка по категории")
	fmt.Println("list - вывести список всех расходов")
	fmt.Println("set - установка бюджета на месяц")
	fmt.Println("export - экспорт в csv")
	fmt.Println("help - вывести список всех команд и их флагов")
	fmt.Println("q - завершить работу приложения")
}
