package domain

type AuthRepository interface {
	CreateUser(user AuthUser) error
	GetUserByUsername(username string) (AuthUser, error)
}
