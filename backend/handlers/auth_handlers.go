package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/leoldding/user-manage/database"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) login(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging In")

	var creds *database.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding JSON body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var storedPass []byte

	err = conn.QueryRow(db.Ctx, "SELECT password FROM users WHERE username = $1;", creds.Username).Scan(&storedPass)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User does not exist: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			log.Printf("Error getting stored password from database: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(storedPass, []byte(creds.Password))
	if err != nil {
		log.Printf("Incorrect password for user: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// create jwt token here

	w.WriteHeader(http.StatusOK)
}

func (db *DB) logout(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging Out")

	// invalidate jwt token here
}
