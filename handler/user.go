package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
	"todo/database"
	"todo/database/dbHelper"
	"todo/models"
	"todo/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserRequest

	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid Method", err)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	v := validator.New()
	if err := v.Struct(user); err != nil {
		http.Error(w, "Failed to validate request body", http.StatusBadRequest)
		return
	}

	// check user exist or not
	exists, existsErr := dbHelper.IsUserExists(user.Email)
	if existsErr != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already  exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user in db
	if saveErr := dbHelper.CreateUser(user.Name, user.Email, hashedPassword); saveErr != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{"User created successfully"})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// login user
func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid Method", err)
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var hashedPassword, userID string
	err := database.DB.QueryRow("SELECT id ,password  FROM users WHERE email = $1", creds.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// generate hash and check it with stored hash
	if !utils.CheckPasswordHash(creds.Password, hashedPassword) {
		err := http.StatusUnauthorized
		http.Error(w, "Unauthorized", err)
		return
	}

	SQL := `INSERT INTO user_session (user_id) VALUES ($1) returning id`
	var sessonID string
	err = database.DB.Get(&sessonID, SQL, userID)
	if err != nil {
		http.Error(w, "Session not inserted in database", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"token": sessonID,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// logout
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid Method", err)
	}
	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		http.Error(w, "Missing session token", http.StatusBadRequest)
		return
	}
	var expiry time.Time
	err := database.DB.QueryRow(` SELECT expiry_at FROM user_session WHERE session_id = $1`, sessionToken).Scan(&expiry)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	_, err = database.DB.Exec(` UPDATE user_session SET expiry_at = $1 WHERE session_id = $2`, time.Now(), sessionToken)
	if err != nil {
		http.Error(w, "Failed to update session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
}
