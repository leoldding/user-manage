package handlers

import (
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging In")
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging Out")
}
