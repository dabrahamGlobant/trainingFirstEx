package storage

import (
	"first-ex/internal/structs"

	"github.com/google/uuid"
)

type Storage interface {
	Create(user structs.User) (structs.User, error)
	Get(userID uuid.UUID) (structs.User, error)
	GetAll() ([]structs.User, error)
	Update(userID uuid.UUID, updatedUser structs.User) (structs.User, error)
	Delete(userID uuid.UUID) error
}
