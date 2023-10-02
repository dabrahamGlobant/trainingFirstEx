package storage

import (
	"first-ex/internal/structs"

	"github.com/google/uuid"
)

type localStorage struct {
	users map[uuid.UUID]structs.User
}

func NewLocalStorage() Storage {
	stu := localStorage{}
	generatedUs := generateUsers()

	stu.users = make(map[uuid.UUID]structs.User)

	for _, user := range generatedUs {
		stu.users[user.ID] = user
	}

	return &stu
}
func (ls *localStorage) Get(uuid uuid.UUID) (structs.User, error) {
	user, ok := ls.users[uuid]
	if !ok {
		return structs.User{}, nil // ERROR SHOULD BE ADDED HERE
	}
	return user, nil

}

func (ls *localStorage) GetAll() ([]structs.User, error) {
	// Pending to add error handling to this case.

	users := make([]structs.User, 0)

	for _, val := range ls.users {
		users = append(users, val)
	}
	return users, nil

}

func (ls *localStorage) Create(user structs.User) (structs.User, error) {
	_, found := ls.users[user.ID]
	if found {
		// Validates if the uuid exists.
		return structs.User{}, nil // ERROR SHOULD BE ADDED HERE
	}
	ls.users[user.ID] = user
	return ls.users[user.ID], nil

}

func (ls *localStorage) Update(uuid uuid.UUID, user structs.User) (structs.User, error) {
	_, found := ls.users[user.ID]
	if !found {
		// Validates if the uuid exists.
		return structs.User{}, nil // ERROR SHOULD BE ADDED HERE
	}
	ls.users[user.ID] = user
	return ls.users[user.ID], nil

}

func (ls *localStorage) Delete(uuid uuid.UUID) error {
	_, found := ls.users[uuid]
	if !found {
		// Validates if the uuid exists.
		return nil // ERROR SHOULD BE ADDED HERE
	}
	delete(ls.users, uuid)
	return nil

}
