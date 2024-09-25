package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (db *DB) getAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Admin getting all users")
}

func (db *DB) updateUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Admin updating user with id %s", id)
}

func (db *DB) deleteUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Admin deleting user with id %s", id)
}
