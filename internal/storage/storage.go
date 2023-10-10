package storage

import (
	"github.com/google/uuid"
)

type Storage interface {
	Create(interface{}) (interface{}, error)
	Get(uuid.UUID) (interface{}, error)
	GetAll() ([]interface{}, error)
	Update(uuid.UUID, interface{}) (interface{}, error)
	Delete(uuid.UUID) error
}
