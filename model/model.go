package model

import "time"



type Expence struct {
	Id     int
	Desc   string
	Amount float64
	Date   time.Time
}
