package user

import (
	"encoding/json"
	"first-ex/internal/storage"
	"first-ex/internal/structs"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	storage storage.Storage //Interfaz storage
}

func NewUserService(storage storage.Storage) UserService {
	return UserService{storage: storage}
}

func (u *UserService) Get(uuid uuid.UUID) (structs.User, error) {

	res, err := u.storage.Get(uuid)
	if err != nil {
		return structs.User{}, structs.ServiceError{
			Code:        structs.ExsistingId,
			Description: err.Error(),
		}
	}
	user, ok := res.(structs.User)
	if !ok {
		user, err = parse(fmt.Sprint(res))
		if err != nil {
			return structs.User{}, structs.ServiceError{
				Code:        structs.Internal,
				Description: structs.ErrJsonParse.Error(),
			}
		}

	}
	return user, nil
}

func (u *UserService) GetAll() ([]structs.User, error) {
	res, err := u.storage.GetAll()
	if err != nil {
		return nil, structs.ServiceError{
			Code:        structs.ConFailed,
			Description: err.Error(),
		}
	}

	users := make([]structs.User, 0)

	for _, v := range res {
		if val, ok := v.(structs.User); ok {
			users = append(users, val)
		} else {
			user, err := parse(fmt.Sprint(v))
			if err != nil {
				return nil, structs.ServiceError{
					Code:        structs.Internal,
					Description: structs.ErrJsonParse.Error(),
				}
			}
			users = append(users, user)
		}
	}
	return users, err
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
	res, err := u.storage.Create(userParsed)
	if err != nil {
		return structs.User{}, structs.ServiceError{
			Code:        structs.ExsistingId,
			Description: err.Error(),
		}
	}
	user, ok := res.(structs.User)
	if !ok {
		user, err = parse(fmt.Sprint(res))
	}
	if err != nil {
		return structs.User{}, structs.ServiceError{
			Code:        structs.Internal,
			Description: structs.ErrJsonParse.Error(),
		}
	}
	return user, err

}

func (u *UserService) Update(uuid uuid.UUID, user structs.User) (structs.User, error) {

	response, err := u.storage.Update(uuid, user)
	if err != nil {
		return structs.User{}, structs.ServiceError{
			Code:        structs.NotFound,
			Description: structs.ErrNotFoundErr.Error(),
		}
	}
	upd, _ := response.(structs.User)

	return upd, nil
}

func (u *UserService) Delete(uuid uuid.UUID) error {

	err := u.storage.Delete(uuid)
	if err != nil {
		return structs.ServiceError{
			Code:        structs.NotFound,
			Description: err.Error(),
		}
	}
	return nil
}

func parse(j string) (structs.User, error) {
	data := structs.User{}
	err := json.Unmarshal([]byte(j), &data)
	if err != nil {
		return structs.User{}, err
	}
	return data, nil
}
