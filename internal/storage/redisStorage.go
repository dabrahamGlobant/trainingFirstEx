package storage

import (
	"context"
	"encoding/json"
	"first-ex/internal/structs"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

/// Interface, similar to User & User Repo.

type RedisStorage struct {
	conn *redis.Client
}

func NewRedisStorage(conn *redis.Client) Storage {
	redis := RedisStorage{conn: conn}

	return &redis
}

func parser(jsonVal string) (structs.User, error) {
	data := structs.User{}
	err := json.Unmarshal([]byte(jsonVal), &data)
	if err != nil {
		return structs.User{}, nil // Pending error handling
	}
	return data, nil
}

func (rs *RedisStorage) Get(uuid uuid.UUID) (structs.User, error) {
	val, errGet := rs.conn.Get(context.Background(), "user_"+uuid.String()).Result()
	if errGet != nil {
		return structs.User{}, nil // Pending error handling
	}
	return parser(val)

}

func (rs *RedisStorage) GetAll() ([]structs.User, error) {
	// Use SCAN or KEYS to fetch all user keys
	keys, err := rs.conn.Keys(context.Background(), "user_*").Result()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store the retrieved users
	var users []structs.User

	// Iterate through the keys and retrieve user data
	for _, key := range keys {
		jsonVal, err := rs.conn.Get(context.Background(), key).Result()
		if err != nil {
			// Handle errors, e.g., skip or log them
			continue
		}

		var user structs.User
		if err := json.Unmarshal([]byte(jsonVal), &user); err != nil {
			// Handle errors, e.g., skip or log them
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (rs *RedisStorage) Create(user structs.User) (structs.User, error) {
	_, found := rs.Get(user.ID)
	if found == nil {
		// Validates if the uuid exists.
		return structs.User{}, nil // ERROR SHOULD BE ADDED HERE
	}
	data, err := json.Marshal(user)

	if err != nil {
		return structs.User{}, err
	}
	err = rs.conn.Set(context.Background(), "user_"+user.ID.String(), data, 0).Err()
	if err != nil {
		panic(err)
	}

	// returns the user
	return rs.Get(user.ID)

}

func (rs *RedisStorage) Update(uuid uuid.UUID, newUser structs.User) (structs.User, error) {
	// Check if user exists
	_, found := rs.Get(uuid)
	if found != nil {
		return structs.User{}, nil
	}
	// Marshal to JSON string
	jsonVal, err := json.Marshal(newUser)
	if err != nil {
		return structs.User{}, err
	}

	// Save in the database
	err = rs.conn.Set(context.Background(), "user_"+newUser.ID.String(), jsonVal, 0).Err()
	if err != nil {
		panic(err)
	}

	// Parse string to User
	return rs.Get(newUser.ID)
}

func (rs *RedisStorage) Delete(uuid uuid.UUID) error {
	// Check if user exists
	_, found := rs.Get(uuid)
	if found != nil {
		return nil // Pending error handler
	}

	// Delete the user from the database
	err := rs.conn.Del(context.Background(), "user_"+uuid.String()).Err()
	if err != nil {
		return err
	}

	return nil
}
