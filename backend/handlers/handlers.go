package handlers

import (
	"context"
	"log"
	"net/http"

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
	// login, consider using this as verify as well
	router.HandleFunc("/login", login).Methods("POST")
	// logout
	router.HandleFunc("/logout", logout).Methods("POST")

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
		// return if not authorized
		next.ServeHTTP(w, r)
	})
}
