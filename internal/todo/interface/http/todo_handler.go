package http

import (
	"context"
	"net/http"
	"todo-app/internal/todo/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type TodoHandler struct {
	uc *usecase.TodoUseCase
}

func NewTodoHandler(api huma.API, uc *usecase.TodoUseCase) {
	handler := &TodoHandler{uc: uc}

	grp := huma.NewGroup(api, "/todos")
	myAuthSecurity := []map[string][]string{
		{"myAuth": {}},
	}
	huma.Register(grp, huma.Operation{
		OperationID: "create-todo",
		Summary:     "Create a new todo item",
		Method:      http.MethodPost,
		Path:        "",
		Security:    myAuthSecurity,
	}, handler.Create)
	huma.Register(grp, huma.Operation{
		OperationID: "list-todos",
		Summary:     "List all todo items",
		Method:      http.MethodGet,
		Path:        "",
		Security:    myAuthSecurity,
	}, handler.List)
	huma.Register(grp, huma.Operation{
		OperationID: "get-todo-by-id",
		Summary:     "Get a todo item by ID",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Security:    myAuthSecurity,
	}, handler.GetByID)
	huma.Register(grp, huma.Operation{
		OperationID: "delete-todo-by-id",
		Summary:     "Delete a todo item by ID",
		Method:      http.MethodDelete,
		Path:        "/{id}",
		Security:    myAuthSecurity,
	}, handler.DeleteByID)
	huma.Register(grp, huma.Operation{
		OperationID: "update-todo-by-id",
		Summary:     "Update a todo item by ID",
		Method:      http.MethodPut,
		Path:        "/{id}",
		Security:    myAuthSecurity,
	}, handler.UpdateByID)
}

func (h *TodoHandler) Create(ctx context.Context, input *CreateTodoInput) (*CreateTodoOutput, error) {
	err := h.uc.CreateTodo(input.Body.Title, input.Body.DueDate)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create todo", err)
	}
	resp := &CreateTodoOutput{}
	resp.Body.Message = "Todo item created successfully"
	return resp, nil
}
func (h *TodoHandler) List(ctx context.Context, input *struct{}) (*ListTodosOutput, error) {
	todos, err := h.uc.GetAllTodos()
	if err != nil {
		return nil, err
	}
	resp := &ListTodosOutput{}
	resp.Body.Todos = todos
	return resp, nil

}
func (h *TodoHandler) GetByID(ctx context.Context, input *struct {
	ID string `path:"id" doc:"ID of the todo item"`
}) (*GetTodoByIDOutput, error) {
	todo, err := h.uc.GetTodoByID(input.ID)
	if err != nil {
		return nil, err
	}
	resp := &GetTodoByIDOutput{}
	resp.Body.Todo = todo
	return resp, nil
}

func (h *TodoHandler) DeleteByID(ctx context.Context, input *struct {
	ID string `path:"id" doc:"ID of the todo item"`
}) (*DeleteTodoOutput, error) {
	err := h.uc.DeleteTodo(input.ID)
	if err != nil {
		return nil, err
	}
	resp := &DeleteTodoOutput{}
	resp.Body.Message = "Todo item deleted successfully"
	return resp, nil
}

func (h *TodoHandler) UpdateByID(ctx context.Context, input *UpdateTodoInput) (*UpdateTodoOutput, error) {
	err := h.uc.UpdateTodo(input.ID, input.Body.Title, input.Body.DueDate)
	if err != nil {
		return nil, err
	}
	resp := &UpdateTodoOutput{}
	resp.Body.Message = "Todo item updated successfully"
	return resp, nil
}
