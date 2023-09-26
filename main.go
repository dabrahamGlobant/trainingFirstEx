package main

import (
	"first-ex/internal/user"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	UserRepository := user.NewUserRepository()

	newUID, _ := uuid.NewRandom()
	newUser := user.User{
		ID:       newUID,
		Name:     "David",
		LastName: "Abraham",
		Email:    "david.abraham@globant.com",
		Active:   true,
		Address: struct {
			City    string
			Country string
			Address string
		}{
			City:    "Cordoba",
			Country: "Argentina",
			Address: "Calle falsa 123",
		},
	}
	UserRepository.Create(newUser)

	// Update an existing user
	updatedUser := user.User{
		ID:       newUID,
		Name:     "Ever",
		LastName: "Cifuentes",
		Email:    "ever.cifuentes@globant.com",
		Active:   true,
		Address: struct {
			City    string
			Country string
			Address string
		}{
			City:    "Cordoba",
			Country: "Argentina",
			Address: "Calle Not Found  404",
		},
	}
	if UserRepository.Update(newUID, updatedUser) {
		fmt.Println("User updated.")
	} else {
		fmt.Println("Failed to update user.")
	}

	// Get a user by ID
	user, found := UserRepository.Get(newUID)
	if found {
		fmt.Println("User found:")
		fmt.Println(user)
	} else {
		fmt.Println("User not found.")
	}

	// Get all users
	allUsers := UserRepository.GetAll()
	fmt.Println("All users:")
	for _, u := range allUsers {
		fmt.Println(u)
	}

	// Delete a user by their ID
	if UserRepository.Delete(newUID) {
		fmt.Println("User deleted.")
	} else {
		fmt.Println("Failed to delete user.")
	}

}
