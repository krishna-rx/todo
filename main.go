package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"todo/database"
	"todo/handler"
	"todo/middlewares"
)

type Login struct {
	HashedPassword string
	SessionToken   string
}

var users = make(map[string]Login)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	fmt.Println(user, pass, host, port, dbname)
	consStr := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable" // getting from compose file
	database.ConnectDB(consStr)                                                                                // connected function to the DB
	defer database.CloseDB()                                                                                   // closing the DB

	// LogIn
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)
	protected.HandleFunc("/logout", handler.Logout).Methods("POST")
	protected.HandleFunc("/CreateTask", handler.CreateTask).Methods("POST")
	protected.HandleFunc("/GetAllTodos", handler.GetAllTodos).Methods("GET")
	protected.HandleFunc("/DeleteTaskById/{id}", handler.DeleteTodo).Methods("DELETE")
	protected.HandleFunc("/UpdateTodoById/{id}", handler.UpdateTodo).Methods("PUT")
	protected.HandleFunc("/", health)

	http.ListenAndServe(":8080", r)
}
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "server running")
}

// user registration

// loging in
