package storage

import (
	"first-ex/internal/structs"

	"github.com/google/uuid"
)

type localStorage struct {
	entities map[uuid.UUID]interface{}
}

func NewLocalStorage() Storage {
	stu := localStorage{}
	generatedUs := GenerateUsers()

	stu.entities = make(map[uuid.UUID]interface{})

	for _, user := range generatedUs {
		stu.entities[user.ID] = user
	}

	return &stu
}
func (ls *localStorage) Get(uuid uuid.UUID) (interface{}, error) {
	entity, ok := ls.entities[uuid]
	if !ok {
		// Generate the new error
		return nil, structs.ErrExistingIdErr
	}
	return entity, nil

}

func (ls *localStorage) GetAll() ([]interface{}, error) {
	// Pending to add error handling to this case.

	entities := make([]interface{}, 0)

	for _, val := range ls.entities {
		entities = append(entities, val)
	}
	return entities, nil

}

func (ls *localStorage) Create(user interface{}) (interface{}, error) {
	uuid := user.(structs.User).ID
	_, found := ls.entities[uuid]
	if found {
		return structs.User{}, structs.ErrExistingIdErr
	}
	ls.entities[uuid] = user
	return ls.entities[uuid], nil

}

func (ls *localStorage) Update(uuid uuid.UUID, user interface{}) (interface{}, error) {
	_, found := ls.entities[uuid]
	if !found {
		return nil, structs.ErrNotFoundErr
	}
	ls.entities[uuid] = user
	return ls.entities[uuid], nil
}

func (ls *localStorage) Delete(uuid uuid.UUID) error {
	_, found := ls.entities[uuid]
	if !found {
		return structs.ErrNotFoundErr
	}
	delete(ls.entities, uuid)
	return nil

}
