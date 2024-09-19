package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) {
	log.Println("Registering Handlers")

	// auth endpoints
	// login, consider using this as verify as well
	router.HandleFunc("/login", login).Methods("POST")
	// logout
	router.HandleFunc("/logout", logout).Methods("POST")

	// user endpoints
	// create user
	router.HandleFunc("/user", createUser).Methods("POST")
	// get user
	router.Handle("/user", authorize(http.HandlerFunc(getUser))).Methods("GET")
	// updated user
	router.Handle("/user", authorize(http.HandlerFunc(updateUser))).Methods("PATCH")
	// delete user
	router.Handle("/user", authorize(http.HandlerFunc(deleteUser))).Methods("DELETE")
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user is authorized
		log.Println("Checking if request is authorized.")
		// return if not authorized
		next.ServeHTTP(w, r)
	})
}
