package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/database"
	"todo/models"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {

	userIDRaw := r.Context().Value("user_id")
	tokeUserID, ok := userIDRaw.(string)
	if !ok {
		http.Error(w, "user_id not found or not a string", http.StatusUnauthorized)
		return
	}
	//userID, err := IsUserSessionValid(tokeUserID)
	//if err != nil {
	//	http.Error(w, "failed to check user authentication", http.StatusInternalServerError)
	//	return
	//}
	//
	//if !userID.Valid {
	//	http.Error(w, "invalid token", http.StatusUnauthorized)
	//	return
	//}
	var todoObj struct {
		Title       string `json:"name"`
		Description string `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&todoObj)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	if todoObj.Title == "" {
		http.Error(w, "Name should is mandatory Invalid JSON", http.StatusBadRequest)
	}
	if todoObj.Description == "" {
		http.Error(w, "Description is mandatory Invalid JSON", http.StatusBadRequest)
	}
	fmt.Println(todoObj.Title, todoObj.Description)
	todoRes := models.TODO{
		Name:        todoObj.Title,
		Description: todoObj.Description,
	}
	_, err = database.DB.Exec(`
  INSERT INTO todos (name, description, user_id)
  VALUES ($1, $2, $3)			
`, todoRes.Name, todoRes.Description, tokeUserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = dbHelper.CreateTodo(user_id, todo.Name, todo.Description)
	//if err != nil {
	//	fmt.Println("not working the todo inserted")
	//	http.Error(w, "can't create the task", http.StatusBadRequest)
	//}

	//w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "task has been created",
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Get all todos
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("user_id")
	fmt.Println("get userID", userIDRaw)
	tokeUserID, ok := userIDRaw.(string)
	if !ok {
		http.Error(w, "user_id not found or not a string", http.StatusUnauthorized)
		return
	}
	//
	//userID, err := IsUserSessionValid(tokeUserID)
	//if err != nil {
	//	http.Error(w, "failed to check user authentication", http.StatusInternalServerError)
	//	return
	//}
	//if !userID.Valid {
	//	http.Error(w, "invalid token", http.StatusUnauthorized)
	//	return
	//}
	var tasks []models.TODO
	query := `SELECT id, name, description FROM todos WHERE user_id = $1`
	var rows *sql.Rows
	rows, err := database.DB.Query(query, tokeUserID)
	fmt.Println(tokeUserID)
	if err != nil {
		http.Error(w, "todos not found", http.StatusNotFound)
		return
	}

	// appending one by one bcz scan only store one type of struct we have array of struct
	for rows.Next() {
		var task models.TODO
		rows.Scan(&task.ID, &task.Name, &task.Description)
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(tasks)
}

// deleted task
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("user_id")
	tokeUserID, ok := userIDRaw.(string)
	if !ok {
		http.Error(w, "user_id not found or not a string", http.StatusUnauthorized)
		return
	}
	///userID, err := IsUserSessionValid(tokeUserID)
	//if err != nil {
	//	http.Error(w, "failed to check user authentication", http.StatusInternalServerError)
	//	return
	//}
	//if !userID.Valid {
	//	http.Error(w, "invalid token", http.StatusUnauthorized)
	//	return
	//}

	vars := mux.Vars(r)
	taskID := vars["id"]
	var userIDTodo string
	err := database.DB.QueryRow(`
	SELECT user_id FROM todos WHERE id = $1
`, taskID).Scan(&userIDTodo)

	if userIDTodo == "" || userIDTodo != tokeUserID {
		http.Error(w, "invalid request to delete todo", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to query todos as user logout", http.StatusInternalServerError)
		return
	}
	_, err = database.DB.Exec("DELETE FROM todos WHERE id = $1", taskID)
	if err != nil {
		http.Error(w, "failed to delete todos", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode([]map[string]string{
		{"message": "task has been deleted"},
	})
}

// updating todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("user_id")
	tokeUserID, ok := userIDRaw.(string)
	if !ok {
		http.Error(w, "user_id not found or not a string", http.StatusUnauthorized)
	}
	//userID, err := IsUserSessionValid(tokeUserID)
	//if err != nil {
	//	http.Error(w, "failed to check user authentication", http.StatusInternalServerError)
	//}
	//if !userID.Valid {
	//	http.Error(w, "invalid token", http.StatusUnauthorized)
	//}

	var req models.UpdateTodo
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	todoId := vars["id"]
	result, err := database.DB.Exec(
		`UPDATE todos SET name = $1, description = $2, is_completed = $3 WHERE id = $4 AND user_id = $5`,
		req.Name, req.Description, req.IsCompleted, todoId, tokeUserID,
	)
	if err != nil {
		http.Error(w, "failed to update todos", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "error fetching affected rows:", http.StatusInternalServerError)
	} else if rowsAffected == 0 {
		log.Println("no todo updated. check if the todo ID and user ID are correct.")
	}
	w.Write([]byte("Todo updated successfully"))
}
