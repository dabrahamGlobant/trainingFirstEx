package storage

import (
	"first-ex/internal/structs"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func generateUsers() []structs.User {

	rand.Seed(time.Now().UnixNano())
	var users []structs.User
	for i := 0; i < 5; i++ {
		user := structs.User{
			ID:       uuid.New(),
			Name:     randomName(),
			LastName: randomName(),
			Email:    randomEmail(),
			Active:   rand.Intn(2) == 0,
			Address: structs.Address{
				City:    randomCity(),
				Country: randomCountry(),
				Address: randomStreet(),
			},
		}
		users = append(users, user)
	}

	return users
}

var firstNames = []string{"Alice", "Bob", "Charlie", "David", "Eva", "Frank", "Grace", "Helen", "Ivy", "Jack"}
var lastNames = []string{"Smith", "Johnson", "Brown", "Lee", "Chen", "Garcia", "Wang", "Kim", "Lopez", "Singh"}
var cities = []string{"New York", "Los Angeles", "Chicago", "San Francisco", "Houston", "Miami", "London", "Paris", "Tokyo", "Sydney"}
var countries = []string{"USA", "Canada", "UK", "France", "Germany", "Japan", "Australia", "India", "Brazil", "China"}

func randomName() string {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func randomEmail() string {
	username := randomName()
	domain := "example.com"
	return fmt.Sprintf("%s@%s", username, domain)
}

func randomCity() string {
	return cities[rand.Intn(len(cities))]
}

func randomCountry() string {
	return countries[rand.Intn(len(countries))]
}

func randomStreet() string {
	return fmt.Sprintf("Random Street %d", rand.Intn(100))
}
