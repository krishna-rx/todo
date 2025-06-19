package handler

import (
	"encoding/json"
	"net/http"
	"todo/models"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")

	// fetch created_at
	IsUserSessionValid(token)

	var todo models.TODO
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid Method", err)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	if todo.Name == "" {
		http.Error(w, "Name should is mandatory Invalid JSON", http.StatusBadRequest)
	}
	if todo.Email == "" {
		http.Error(w, "Email is mandatory Invalid JSON", http.StatusBadRequest)
	}
	if todo.Description == "" {
		http.Error(w, "Description is mandatory Invalid JSON", http.StatusBadRequest)
	}

}
