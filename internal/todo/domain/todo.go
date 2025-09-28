package domain

import "time"

type Todo struct {
	ID      string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Unique identifier for the todo item"`
	Title   string    `json:"title" example:"Buy milk" doc:"Title of the todo item"`
	DueDate time.Time `json:"dueDate" example:"2023-10-10T10:00:00Z" doc:"Due date of the todo item"`
	Done    bool      `json:"done" example:"false" doc:"Completion status of the todo item"`
}
