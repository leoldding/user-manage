package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/leoldding/user-manage/database"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating User")

	// decode request body into user
	var newUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Printf("Error decoding JSON body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get database connection from pool
	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// start database transaction
	tx, err := conn.Begin(db.Ctx)
	if err != nil {
		log.Printf("Error beginning database transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(db.Ctx)

	// TODO: check if the username already exists

	// use bcrypt to hash the users password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		log.Printf("Error hashing user password: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert data users table
	_, err = tx.Exec(db.Ctx, "INSERT INTO users (username, password, first_name, last_name) VALUES ($1, $2, $3, $4);", newUser.Username, hashedPass, newUser.FirstName, newUser.LastName)
	if err != nil {
		log.Printf("Error inserting users table: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert data into user_roles table
	_, err = tx.Exec(db.Ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ((SELECT id FROM users WHERE username = $1), 2);", newUser.Username)
	if err != nil {
		log.Printf("Error inserting into user_roles table: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// commit transaction
	err = tx.Commit(db.Ctx)
	if err != nil {
		log.Printf("Error commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (db *DB) getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting User")

	claims := r.Context().Value("userClaims").(jwt.MapClaims)
	id := claims["id"]

	// get database connection from pool
	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// retrieve user information from database
	var user *database.User

	err = conn.QueryRow(db.Ctx, "SELECT username, first_name, last_name FROM users WHERE id = $1;", id).Scan(&user)
	if err != nil {
		log.Printf("Error retrieving user from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (db *DB) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating User")

	claims := r.Context().Value("userClaims").(jwt.MapClaims)
	id := claims["id"]

	// decode request body into user
	var updateUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		log.Printf("Error decoding JSON body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get database connection from pool
	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// TODO: should separate updating each field, transaction

	// update user values
	_, err = conn.Exec(db.Ctx, "UPDATE users SET username = $2, password = $3, first_name = $4, last_name = $5 WHERE id = $1;", id, updateUser.Username, updateUser.Password, updateUser.FirstName, updateUser.LastName)
	if err != nil {
		log.Printf("Error updating user values: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *DB) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting User")

	claims := r.Context().Value("userClaims").(jwt.MapClaims)
	id := claims["id"]

	// decode request body into user
	var deleteUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&deleteUser); err != nil {
		log.Printf("Error decoding JSON body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get database connection from pool
	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// start database transaction
	tx, err := conn.Begin(db.Ctx)
	if err != nil {
		log.Printf("Error beginning database transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(db.Ctx)

	// delete data from user_roles table
	_, err = tx.Exec(db.Ctx, "DELETE FROM user_roles WHERE user_id = $1;", id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete data from users table
	_, err = tx.Exec(db.Ctx, "DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// commit transaction
	err = tx.Commit(db.Ctx)
	if err != nil {
		log.Printf("Error commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
