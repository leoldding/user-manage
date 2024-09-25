package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user is authorized
		log.Println("Checking if request is authorized.")

		// retrieve jwt
		cookie, err := r.Cookie("user-jwt")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("Missing JWT: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				log.Printf("Invalid JWT: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		// parse token
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.Printf("Incorrect JWT signing method: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return nil, errors.New("Incorrect signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Printf("Error parsing JWT: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// invalid token
		if !token.Valid {
			log.Printf("Invalid JWT: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Unable to extract claims: %v", ok)
			http.Error(w, "No claims found", http.StatusInternalServerError)
			return
		}

		log.Println(claims["user"])
		ctx := context.WithValue(r.Context(), "userClaims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("USER AUTHORIZATION")
		next.ServeHTTP(w, r)
	})
}

func adminAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("ADMIN AUTHORIZATION")
		next.ServeHTTP(w, r)
	})
}
