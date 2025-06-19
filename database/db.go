package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var DB *sqlx.DB
var err error

// Get, Select, Exec

func ConnectDB(constr string) {
	fmt.Println(constr)
	DB, err = sqlx.Open("postgres", constr) // connecting to the database using sql prebuild
	if err != nil {
		logrus.Errorf("failed to connect to postgres: %v", err)
	}

	err = DB.Ping() // check if it's running or not
	if err != nil {
		logrus.Errorf("failed to ping postgres: %v", err)
	} else {
		fmt.Println("Successfully connected to postgres")
	}
	CreateMigrationsTable(DB)
}
func CloseDB() {
	err := DB.Close()
	if err != nil {
		logrus.Errorf("failed to close postgres: %v", err)
	}
}
