package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	LastName string
	Email    string
	Active   bool
	Address  Address
}

type Address struct {
	City    string
	Country string
	Address string
}
