package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func RegisterHandlers(router *mux.Router, pool *pgxpool.Pool, ctx context.Context) {
	log.Println("Registering Handlers")

	db := DB{pool, ctx}

	// auth endpoints
	// login
	router.HandleFunc("/login", db.login).Methods("POST")
	// logout
	router.HandleFunc("/logout", db.logout).Methods("GET")
	// authenticate
	router.HandleFunc("/auth", isAuthenticated).Methods("GET")

	// user endpoints
	// create user
	router.HandleFunc("/user", db.createUser).Methods("POST")
	// get user
	router.Handle("/user", authorize(http.HandlerFunc(db.getUser))).Methods("GET")
	// updated user
	router.Handle("/user", authorize(http.HandlerFunc(db.updateUser))).Methods("PATCH")
	// delete user
	router.Handle("/user", authorize(http.HandlerFunc(db.deleteUser))).Methods("DELETE")
}

func authorize(next http.Handler) http.Handler {
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
