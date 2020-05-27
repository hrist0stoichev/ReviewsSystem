package models

type User struct {
	Id             string
	Email          string
	EmailConfirmed bool
	HashedPassword string
	Role           Role
}
