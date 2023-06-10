package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

const Secret = "jwt-secret>"

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "1Unauthenticated User", http.StatusBadRequest)
			default:
				http.Error(w, "2Unauthenticated User", http.StatusInternalServerError)
			}
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "3Unauthenticated User", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserId(w http.ResponseWriter, r *http.Request) (uint, error) {
	cookie, _ := r.Cookie("jwt")

	token, _ := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})

	payload := token.Claims.(*jwt.StandardClaims)

	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}
