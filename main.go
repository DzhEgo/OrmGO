package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"testOrmGo/Users"
)

func startServer() {
	router := mux.NewRouter()

	router.HandleFunc("/users", Users.CreateUser).Methods("POST")
	router.HandleFunc("/users", Users.AllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", Users.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", Users.UpdateUser).Methods("PUT")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	dsn := "user=postgres password=Lax212212 dbname=testUser sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Users.Users{})

	fmt.Println("Запуск сервера...")
	startServer()
}
