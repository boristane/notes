package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// UserClaims is the data that is inside the authentication token
type UserClaims struct {
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	ID       uint64 `json:"id"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_USERS_KEY"))

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		bearerToken := strings.Split(tokenString, " ")
		if tokenString == "" || len(bearerToken) != 2 {
			send401(w)
			return
		}
		claims := &UserClaims{}
		tkn, err := jwt.
			ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil || !tkn.Valid {
			send401(w)
			return
		}

		claims, ok := tkn.Claims.(*UserClaims)

		if !ok {
			send401(w)
			return
		}

		ctx := context.WithValue(r.Context(), user_id, claims.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type key int

const user_id key = iota
