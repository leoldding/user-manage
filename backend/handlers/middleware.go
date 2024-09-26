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
		// check if user is authenticated
		log.Println("Checking if user is authenticated")

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

func adminAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Admin authorization")

		// retrieve claims
		claims := r.Context().Value("userClaims").(jwt.MapClaims)
		role := claims["role"]

		// check if claims has admin role
		if role != "admin" {
			log.Println("Not an authorized admin")
			http.Error(w, "Unauthorized admin", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
