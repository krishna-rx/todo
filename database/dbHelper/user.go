package dbHelper

import "todo/database"

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT count(id) > 0 as is_exist
            FROM users
            WHERE email = TRIM($1)
              AND archived_at IS NULL`

	var check bool
	chkErr := database.DB.Get(&check, SQL, email)
	return check, chkErr
}
func CreateUser(name, email, password string) error {
	SQL := `INSERT INTO users (name, email, password)
            VALUES (TRIM($1), TRIM($2), $3)`

	_, crtErr := database.DB.Exec(SQL, name, email, password)
	return crtErr
}
