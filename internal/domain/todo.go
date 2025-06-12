package domain

import "time"

type Todo struct {
	ID      string
	Title   string
	DueDate time.Time
}
