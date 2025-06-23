package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "invalid token format dont found bearer", http.StatusUnauthorized)
			return
		}

		//getting parse check the algorithm
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// get the user ID and provide it in context
		rawUserID, exists := claims["user_id"]

		if !exists || rawUserID == nil {
			http.Error(w, "user_id claim missing", http.StatusUnauthorized)
			return
		}

		userID, ok := rawUserID.(string)
		fmt.Println(ok, rawUserID, userID)
		if !ok {
			http.Error(w, "user_id claim is not a string", http.StatusUnauthorized)
			return
		}

		// Now userID is safe to use
		fmt.Println(userID)
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
