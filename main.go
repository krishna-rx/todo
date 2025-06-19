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
	r := mux.NewRouter()
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/logout", handler.Logout).Methods("POST")
	r.HandleFunc("/CreateTask", handler.CreateTask).Methods("POST")
	r.HandleFunc("/", health)

	http.ListenAndServe(":8080", r)
}
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "server running")
}

// user registration

// loging in
