package domain

type TokenGenerator interface {
	Generate(user AuthUser) (string, error)
}
