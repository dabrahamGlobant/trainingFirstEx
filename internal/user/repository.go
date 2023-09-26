package user

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRepository struct {
	users map[uuid.UUID]User
}

func NewUserRepository() *UserRepository {

	users := make(map[uuid.UUID]User)

	for i := 0; i < 5; i++ {
		id, _ := uuid.NewRandom()
		user := User{
			ID:       id,
			Name:     fmt.Sprintf("User%d", i+1),
			LastName: fmt.Sprintf("Lastname%d", i+1),
			Email:    fmt.Sprintf("user%d@gmail.com", i+1),
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
		users[id] = user
	}
	return &UserRepository{users}

}

func (response *UserRepository) Create(user User) {
	response.users[user.ID] = user
}

func (response *UserRepository) Get(userID uuid.UUID) (User, bool) {
	user, ok := response.users[userID]
	return user, ok
}

func (response *UserRepository) GetAll() map[uuid.UUID]User {
	return response.users
}

func (response *UserRepository) Update(userID uuid.UUID, updatedUser User) bool {
	_, ok := response.users[userID]
	if !ok {
		fmt.Print("User not found")
		return false
	}
	response.users[userID] = updatedUser
	return true
}

func (response *UserRepository) Delete(userID uuid.UUID) bool {
	_, ok := response.users[userID]
	if !ok {
		fmt.Print("User not found")
		return false
	}
	delete(response.users, userID)
	return true
}
