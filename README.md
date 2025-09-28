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

Required variables:

- MONGO_DB_URI: MongoDB connection string
- MONGO_DB_NAME: Database name
- JWT_SECRET: Secret for signing JWT tokens
- SERVER_ADDRESS (optional): Defaults to localhost:8080
- AUTH_REPO (optional): memory (default) or mongo

To persist auth users to MongoDB, set `AUTH_REPO=mongo`. Users will be stored in the `auth_users` collection with a unique index on the `username` field.

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

