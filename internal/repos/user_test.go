package user_test

import (
	"errors"
	user "first-ex/internal/repos"
	"first-ex/internal/storage"
	"first-ex/internal/structs"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type mockStorage struct {
	data map[uuid.UUID]interface{}
}

func (ms *mockStorage) Get(uuid uuid.UUID) (interface{}, error) {
	if val, ok := ms.data[uuid]; ok {
		return val, nil
	}
	return nil, errors.New("not found")
}

func (ms *mockStorage) GetAll() ([]interface{}, error) {
	var entities []interface{}
	for _, val := range ms.data {
		entities = append(entities, val)
	}
	return entities, nil
}

func (ms *mockStorage) Create(user interface{}) (interface{}, error) {
	return nil, errors.New("create error")
}

func (ms *mockStorage) Update(uuid uuid.UUID, user interface{}) (interface{}, error) {
	return nil, errors.New("update error")
}

func (ms *mockStorage) Delete(uuid uuid.UUID) error {
	return errors.New("delete error")
}

func TestUserService_Get(t *testing.T) {
	mockData := make(map[uuid.UUID]interface{})
	users := storage.GenerateUsers()

	for _, user := range users {
		mockData[user.ID] = user
	}
	mockStorage := &mockStorage{data: mockData}
	userService := user.NewUserService(mockStorage)

	t.Run("TestGetExistingUser", func(t *testing.T) {
		user, err := userService.Get(users[0].ID)
		if !errors.Is(err, nil) {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		if !reflect.DeepEqual(user.ID, users[0].ID) {
			t.Errorf("Expected user with ID %v, but got %v", users[0].ID, user.ID)
		}
	})

	t.Run("TestGetNonExistingUser", func(t *testing.T) {
		nonExistingUUID := uuid.New()
		expectedError := structs.ServiceError{
			Code:        structs.NotFound,
			Description: structs.ErrNotFoundErr.Error(),
		}
		_, err := userService.Get(nonExistingUUID)
		if !errors.Is(err, expectedError) {
			t.Errorf("Expected an error, but got no error.")
		}
	})
}

func TestUserService_GetAll(t *testing.T) {
	mockData := make(map[uuid.UUID]interface{})
	users := storage.GenerateUsers()
	for _, user := range users {
		mockData[user.ID] = user
	}
	mockStorage := &mockStorage{data: mockData}
	userService := user.NewUserService(mockStorage)

	t.Run("TestGetAllUsers", func(t *testing.T) {
		users, err := userService.GetAll()
		if !errors.Is(err, nil) {
			t.Errorf("Expected no error, but got an error: %v", err)
		}
		if !reflect.DeepEqual(len(users), len(mockData)) {
			t.Errorf("Expected %d users, but got %d", len(mockData), len(users))
		}
	})
}

func TestUserService_Delete(t *testing.T) {
	mockData := make(map[uuid.UUID]interface{})
	users := storage.GenerateUsers()
	for _, user := range users {
		mockData[user.ID] = user
	}
	mockStorage := &mockStorage{data: mockData}
	userService := user.NewUserService(mockStorage)

	t.Run("TestDeleteNonExistingUser", func(t *testing.T) {
		nonExistingUUID := uuid.New()

		err := userService.Delete(nonExistingUUID)
		expectedError := structs.ServiceError{
			Code:        structs.NotFound,
			Description: structs.ErrNotFoundErr.Error(),
		}

		if !errors.Is(err, expectedError) {
			t.Errorf("Expected an error, but got no error.")
		}
	})
}
func TestUserService_Update(t *testing.T) {
	mockData := make(map[uuid.UUID]interface{})
	users := storage.GenerateUsers()
	for _, user := range users {
		mockData[user.ID] = user
	}
	mockStorage := &mockStorage{data: mockData}
	userService := user.NewUserService(mockStorage)

	t.Run("TestUpdateNonExistingUser", func(t *testing.T) {

		nonExistingUUID := uuid.New()
		modifiedUser := structs.User{}

		_, err := userService.Update(nonExistingUUID, modifiedUser)
		expectedError := structs.ServiceError{
			Code:        structs.NotFound,
			Description: structs.ErrNotFoundErr.Error(),
		}

		if !errors.Is(err, expectedError) {
			t.Errorf("Expected an error, but got no error.")
		}
	})
}
func TestParse(t *testing.T) {
	t.Run("TestJsonParsing", func(t *testing.T) {
		jsonStr := `{"ID":"123e4567-e89b-12d3-a456-426655440000","Name":"John Doe","LastName":"Doe","Email":"john.doe@example.com","Active":true}`

		parsedUser, err := user.Parse(jsonStr)
		if err != nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Verify the parsed user's fields.
		expectedUser := structs.User{
			ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"),
			Name:     "John Doe",
			LastName: "Doe",
			Email:    "john.doe@example.com",
			Active:   true,
		}

		if !reflect.DeepEqual(parsedUser, expectedUser) {
			t.Errorf("Parsed user does not match the expected user.")
		}

	})
	t.Run("WrongJsonParsing", func(t *testing.T) {
		invalidJSONStr := `{"ID": "123" "Name": "John Doe"}`

		parsedUser, err := user.Parse(invalidJSONStr)
		if err == nil {
			t.Errorf("Expected no error, but got an error: %v", err)
		}

		// Verify the parsed user's fields.
		expectedUser := structs.User{
			ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"),
			Name:     "John Doe",
			LastName: "Doe",
			Email:    "john.doe@example.com",
			Active:   true,
		}

		if reflect.DeepEqual(parsedUser, expectedUser) {
			t.Errorf("Parsed user does not match the expected user.")
		}

	})

}

// func TestUserService_Create(t *testing.T) {
// 	mockData := make(map[uuid.UUID]interface{})
// 	users := storage.GenerateUsers()
// 	for _, user := range users {
// 		mockData[user.ID] = user
// 	}
// 	mockStorage := &mockStorage{data: mockData}
// 	userService := user.NewUserService(mockStorage)

// 	t.Run("TestCreateUser", func(t *testing.T) {
// 		newUserRequest := structs.UserRequest{
// 			Name:     "New User",
// 			LastName: "Last Name",
// 			Email:    "newuser@example.com",
// 			Active:   true,
// 			Address: structs.AddressRequest{
// 				City:    "Fake City",
// 				Country: "Inflational Argentina",
// 				Address: "Fake Street 1234",
// 			},
// 		}

// 		newUser, err := userService.Create(newUserRequest)
// 		if err != nil {
// 			t.Errorf("Expected no error, but got an error: %v", err)
// 		}

// 		if newUser.Name != newUserRequest.Name {
// 			t.Errorf("Expected user's Name to be %s, but got %s", newUserRequest.Name, newUser.Name)
// 		}
// 		if newUser.Email != newUserRequest.Email {
// 			t.Errorf("Expected user's Email to be %s, but got %s", newUserRequest.Email, newUser.Email)
// 		}

// 	})
// }
