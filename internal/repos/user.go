package user

import (
	"first-ex/internal/storage"
	"first-ex/internal/structs"

	"github.com/google/uuid"
)

type UserService struct {
	storage storage.Storage //Interfaz storage
}

func NewUserService(storage storage.Storage) UserService {
	return UserService{storage: storage}
}

func (u *UserService) Get(uuid uuid.UUID) (structs.User, error) {
	return u.storage.Get(uuid)
}

func (u *UserService) GetAll() ([]structs.User, error) {
	return u.storage.GetAll()
}
func (u *UserService) Create(req structs.UserRequest) (structs.User, error) {
	newUuid := uuid.New()
	userParsed := structs.User{
		ID:       newUuid,
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.LastName,
		Active:   req.Active,
		Address:  structs.Address(req.Address),
	}
	user, err := u.storage.Create(userParsed)
	return user, err
}

func (u *UserService) Update(uuid uuid.UUID, user structs.User) (structs.User, error) {

	return u.storage.Update(uuid, user)
}

func (u *UserService) Delete(uuid uuid.UUID) error {
	return u.storage.Delete(uuid)
}
