package domain

type AuthUser struct {
	Username     string
	PasswordHash string
	Token        string
}
