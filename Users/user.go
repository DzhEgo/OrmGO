package Users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type Users struct {
	gorm.Model
	id       uint   `gorm:"primaryKey"`
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

	err := db.First(&user, id)
	if err.Error != nil {
		http.Error(w, "Пользователь не найден!", http.StatusNotFound)
		return
	}

	err = db.Where("id_user = ?", id).Delete(&user)
	if err.Error != nil {
		http.Error(w, "Ошибка при удалении пользователя!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Пользователь удален!")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user Users

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	search := db.First(&user, id)
	if search.Error != nil {
		http.Error(w, "Пользователь не найден!", http.StatusNotFound)
		return
	}

	res := json.Unmarshal(body, &user)
	if res != nil {
		http.Error(w, "Ошибка при декодировании JSON!", http.StatusBadRequest)
		return
	}

	upd := db.Where("id_user = ?", id).Save(&user)
	if upd.Error != nil {
		http.Error(w, "Ошибка при обновлении пользователя!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Пользователь обновлен!")
}
