package http

import "todo-app/internal/domain"

type (
	CreateTodoInput struct {
		Body struct {
			Title   string `json:"title" doc:"Title of the todo item" example:"Buy groceries"`
			DueDate string `json:"due_date" doc:"Due date for the todo item" example:"2023-10-01"`
		}
	}
	UpdateTodoInput struct {
		ID   string `path:"id" doc:"ID of the todo item"`
		Body struct {
			Title   string `json:"title" doc:"Title of the todo item" example:"Buy groceries"`
			DueDate string `json:"due_date" doc:"Due date for the todo item" example:"2023-10-01"`
		}
	}
)

type (
	CreateTodoOutput struct {
		Body struct {
			Message string `json:"message" example:"Todo item created successfully" doc:"Confirmation message"`
		}
	}
	ListTodosOutput struct {
		Body struct {
			Todos []*domain.Todo `json:"todos" doc:"List of todo items"`
		}
	}

	GetTodoByIDOutput struct {
		Body struct {
			Todo *domain.Todo `json:"todo" doc:"Todo item"`
		}
	}

	DeleteTodoOutput struct {
		Body struct {
			Message string `json:"message" example:"Todo item deleted successfully" doc:"Confirmation message"`
		}
	}
	UpdateTodoOutput struct {
		Body struct {
			Message string `json:"message" example:"Todo item updated successfully" doc:"Confirmation message"`
		}
	}
)
