package model

import "time"

type Expense struct {
	Id     int
	Desc   string
	Amount float64
	Date   time.Time
}

type ExpenseUpdate struct {
	ID          int
	Description *string
	Amount      *float64
}
