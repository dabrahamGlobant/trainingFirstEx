package structs

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

type UserRequest struct {
	Name     string `json:"name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Active   bool
	Address  AddressRequest `json:"address" validate:"required"`
}

type AddressRequest struct {
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
	Address string `json:"address" validate:"required"`
}
