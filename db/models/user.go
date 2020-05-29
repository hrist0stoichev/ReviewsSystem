package models

type User struct {
	Id                     string
	Email                  string
	EmailConfirmed         bool
	EmailConfirmationToken *string
	HashedPassword         string
	Role                   Role
}
