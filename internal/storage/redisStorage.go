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
	conn   *redis.Client
	prefix string
}

func NewRedisStorage(conn *redis.Client) Storage {
	redis := RedisStorage{conn: conn, prefix: "user_"}

	return &redis
}

func (rs *RedisStorage) Get(uuid uuid.UUID) (interface{}, error) {
	val, errGet := rs.conn.Get(context.Background(), rs.prefix+uuid.String()).Result()
	if errGet != nil {
		return structs.User{}, errGet
	}
	return val, nil

}

func (rs *RedisStorage) GetAll() ([]interface{}, error) {
	// Use SCAN or KEYS to fetch all user keys
	keys, err := rs.conn.Keys(context.Background(), rs.prefix+"*").Result()
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store the retrieved users
	var users []interface{}

	// Iterate through the keys and retrieve user data
	for _, key := range keys {
		jsonVal, err := rs.conn.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		var user structs.User
		if err := json.Unmarshal([]byte(jsonVal), &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (rs *RedisStorage) Create(entity interface{}) (interface{}, error) {
	uuid := entity.(structs.User).ID
	_, found := rs.Get(uuid)
	if found == nil {
		// Validates if the uuid exists.
		return structs.User{}, found
	}
	data, err := json.Marshal(entity)

	if err != nil {
		return structs.User{}, err
	}
	err = rs.conn.Set(context.Background(), rs.prefix+uuid.String(), data, 0).Err()
	if err != nil {
		panic(err)
	}

	// returns the user
	return rs.Get(uuid)

}

func (rs *RedisStorage) Update(uuid uuid.UUID, entity interface{}) (interface{}, error) {
	// Check if user exists
	_, found := rs.Get(uuid)
	if found != nil {
		return structs.User{}, structs.ErrNotFoundErr
	}
	// Marshal to JSON string
	jsonVal, err := json.Marshal(entity)
	if err != nil {
		return structs.User{}, structs.ErrJsonParse
	}

	// Save in the database
	err = rs.conn.Set(context.Background(), "user_"+uuid.String(), jsonVal, 0).Err()
	if err != nil {
		return nil, err
	}

	// Parse string to User
	return rs.Get(uuid)
}

func (rs *RedisStorage) Delete(uuid uuid.UUID) error {
	// Check if user exists
	_, found := rs.Get(uuid)
	if found != nil {
		return structs.ErrExistingIdErr
	}

	// Delete the user from the database
	err := rs.conn.Del(context.Background(), rs.prefix+uuid.String()).Err()
	if err != nil {
		return err
	}

	return nil
}
