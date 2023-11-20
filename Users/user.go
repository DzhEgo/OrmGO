package Users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type Users struct {
	ID_User  int    `json:"ID_User"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var db *gorm.DB

func init() {
	var err error
	dsn := "user=postgres password=Lax212212 dbname=testUser sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Users{})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user Users

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := db.Create(&user)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Пользователь создан!")
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	var users []Users
	result := db.Find(&users)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user Users

	res := db.First(&user, id)

	fmt.Println(user)
	fmt.Println(id)
	if res.Error != nil {
		http.Error(w, "Пользователь не найден!", http.StatusNotFound)
		return
	}

	if err := db.Delete(&user); err != nil {
		http.Error(w, "Ошибка при удалении!", http.StatusInternalServerError)
		return
	}
	fmt.Println(user)

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user Users

	db.First(&user, id)
	json.NewDecoder(r.Body).Decode(&user)
	db.Model(&user).Updates(&user)
	json.NewEncoder(w).Encode(user)
}
