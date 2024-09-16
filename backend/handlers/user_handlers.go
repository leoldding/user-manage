package handlers

import (
	"log"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating User")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating User")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting User")
}
