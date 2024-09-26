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
	// login
	router.HandleFunc("/login", db.login).Methods("POST")
	// logout
	router.HandleFunc("/logout", db.logout).Methods("GET")
	// authenticate
	router.Handle("/auth", authenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))).Methods("GET")

	// user endpoints
	// create user
	router.HandleFunc("/user", db.createUser).Methods("POST")
	// get user
	router.Handle("/user", authenticationMiddleware(http.HandlerFunc(db.getUser))).Methods("GET")
	// update user
	router.Handle("/user", authenticationMiddleware(http.HandlerFunc(db.updateUser))).Methods("PUT")
	// delete user
	router.Handle("/user", authenticationMiddleware(http.HandlerFunc(db.deleteUser))).Methods("DELETE")

	// admin endpoints
	// get users
	router.Handle("/users", authenticationMiddleware(adminAuthorizationMiddleware(http.HandlerFunc(db.getAllUsers)))).Methods("GET")
	// update user
	router.Handle("/user/{id}", authenticationMiddleware(adminAuthorizationMiddleware(http.HandlerFunc(db.updateUserById)))).Methods("PUT")
	// delete user
	router.Handle("/user/{id}", authenticationMiddleware(adminAuthorizationMiddleware(http.HandlerFunc(db.deleteUserById)))).Methods("DELETE")
}
