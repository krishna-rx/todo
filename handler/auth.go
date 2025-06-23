package handler

import (
	"fmt"
	"github.com/volatiletech/null"
	"time"
	"todo/database"
	"todo/models"
)

func IsUserSessionValid(token string) (null.String, error) {
	// Struct with nullable ArchivedAt

	var sessionData models.SessionData

	// Fetch session data by session ID (token)
	err := database.DB.Get(&sessionData, `SELECT user_id, created_at, archived_at FROM user_session WHERE id = $1`, token)
	if err != nil {
		fmt.Printf("DB error: %+v\n", err)
		return sessionData.UserID, err
	}

	now := time.Now().UTC()
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)

	// Validate session:
	// - Created within the last 7 days
	// - Not archived (archived_at is NULL or in the future)
	if sessionData.CreatedAt.Before(sevenDaysAgo) {
		fmt.Println("session too old.")
		return sessionData.UserID, nil
	}

	if sessionData.ArchivedAt != nil && sessionData.ArchivedAt.Before(now) {
		fmt.Println("session archived.")
		return sessionData.UserID, nil
	}

	return sessionData.UserID, nil
}
