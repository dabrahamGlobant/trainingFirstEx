package main

import (
	"first-ex/internal/handlers"
	user "first-ex/internal/repos"
	"first-ex/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	router := mux.NewRouter()
	var usersService user.UserService
	//Read .env to choose if we should use localstorage or redis
	switch goDotEnvVariable("STORAGE") {
	case "Redis":
		fmt.Println("REDIS_URL")
		var connection = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf(goDotEnvVariable("REDIS_URL")),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		usersService = user.NewUserService(storage.NewRedisStorage(connection))
	default:
		fmt.Println("LOCAL STORAGE")
		usersService = user.NewUserService(storage.NewLocalStorage())
	}
	//declare the user service injecting the storage dependency

	//ROUTER
	router.HandleFunc("/users", handlers.GetAll(usersService)).Methods("GET")
	router.HandleFunc("/users", handlers.Post(usersService)).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.Get(usersService)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.Delete(usersService)).Methods("DELETE")
	router.HandleFunc("/users/{id}", handlers.Put(usersService)).Methods("PUT")

	fmt.Println("Running on PORT 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// Auxiliar function to get a .env var
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
