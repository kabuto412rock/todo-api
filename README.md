Learning Clean Architecture from ChatGPT and implementing it with Golang
## Description

A simple Todo REST API built with Go, following Clean Architecture principles.

## Prerequisites

- Go 1.18 or higher
- MongoDB instance
- [Optional] Docker for running MongoDB locally

## Environment Variables
Create a `.env` file at the project root and set the following variables:
```bash
cp .env.example .env
```

## Installation
```bash
go mod tidy
```

## Running the server

```bash
make run
```

## API Endpoints

| Method | Route      | Description             |
| ------ | ---------- | ----------------------- |
| GET    | /todos     | List all todos          |
| GET    | /todos/:id | Get a single todo       |
| POST   | /todos     | Create a new todo       |
| PUT    | /todos/:id | Update an existing todo |
| DELETE | /todos/:id | Delete a todo           |

### Auth Endpoints

| Method | Route        | Description          |
| ------ | ------------ | -------------------- |
| POST   | /auth/login  | User login           |
| POST   | /auth/register | User registration    |
## Running Tests

```bash
make test
```

