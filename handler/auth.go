package handler

import (
	"time"
	"todo/database"
)

func IsUserSessionValid(token string) (bool, error) {
	// sessionId se created_at and archived_at
	var created, archived time.Time
	err := database.DB.QueryRow(` SELECT expiry_at FROM user_session WHERE session_id = $1`, token).Scan(&created, &archived)
	if err != nil {
		return false, err
	}
	check := time.Now().Add(7 * 24 * time.Hour).UTC()
	if !created.Equal(check) && archived == time.Time {
		return false, nil
	}
	return true, nil
	// time.Now().Add(7 * 24 * time.Hour)  use UTC()
	// if crated_at + 7 days < curent_time && arcived_at is null
}
