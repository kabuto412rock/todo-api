package http

import (
	"time"
	"todo-app/internal/todo/domain"
)

type (
	ListQueryParams struct {
		Page  int `query:"page" doc:"Page number for pagination" example:"0"`
		Limit int `query:"limit" doc:"Number of items per page" example:"10"`
	}
	ListTodosInput struct {
		ListQueryParams
		Title string `query:"title" doc:"Filter todos by title" example:"groceries"`
	}
	CreateTodoInput struct {
		Body struct {
			Title   string    `json:"title" doc:"Title of the todo item" example:"Buy groceries"`
			DueDate time.Time `json:"dueDate" doc:"Due date for the todo item" example:"2023-10-10T10:00:00Z"`
			Done    bool      `json:"done" doc:"Completion status of the todo item" example:"false"`
		}
	}
	UpdateTodoInput struct {
		ID   string `path:"id" doc:"ID of the todo item"`
		Body struct {
			Title   string    `json:"title" doc:"Title of the todo item" example:"Buy groceries"`
			DueDate time.Time `json:"dueDate" doc:"Due date for the todo item" example:"2023-10-10T10:00:00Z"`
			Done    bool      `json:"done" doc:"Completion status of the todo item" example:"false"`
		}
	}
)

type (
	ListResponseMeta struct {
		Page  int   `json:"page" example:"0" doc:"Current page number"`
		Limit int   `json:"limit" example:"10" doc:"Number of items per page"`
		Total int64 `json:"total" example:"100" doc:"Total number of items"`
	}
	CreateTodoOutput struct {
		Body struct {
			Message string `json:"message" example:"Todo item created successfully" doc:"Confirmation message"`
		}
	}
	ListTodosOutput struct {
		Body struct {
			Data []*domain.Todo   `json:"data" doc:"List of todo items"`
			Meta ListResponseMeta `json:"meta" doc:"Pagination metadata"`
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
