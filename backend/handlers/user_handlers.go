package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leoldding/user-manage/database"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating User")

	var newUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Printf("Error decoding JSON body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		log.Printf("Error hashing user password: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(db.Ctx)
	if err != nil {
		log.Printf("Error beginning database transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(db.Ctx)

	_, err = tx.Exec(db.Ctx, "INSERT INTO users (username, password, first_name, last_name) VALUES ($1, $2, $3, $4);", newUser.Username, hashedPass, newUser.FirstName, newUser.LastName)
	if err != nil {
		log.Printf("Error inserting users table: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(db.Ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ((SELECT id FROM users WHERE username = $1), 2);", newUser.Username)
	if err != nil {
		log.Printf("Error inserting into user_roles table: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(db.Ctx, "SELECT id, username, first_name, last_name FROM users;")
	if err != nil {
		log.Printf("Error retrieving users from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []*database.User
	for rows.Next() {
		var user database.User
		err := rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName)
		if err != nil {
			log.Printf("Error marshaling database values into variables: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, &user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (db *DB) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating User")

	var updateUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
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

	_, err = conn.Exec(db.Ctx, "UPDATE users SET username = $2, password = $3, first_name = $4, last_name = $5 WHERE id = $1;", updateUser.Id, updateUser.Username, updateUser.Password, updateUser.FirstName, updateUser.LastName)
	if err != nil {
		log.Printf("Error updating user values: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *DB) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting User")

	var deleteUser *database.User
	if err := json.NewDecoder(r.Body).Decode(&deleteUser); err != nil {
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

	tx, err := conn.Begin(db.Ctx)
	if err != nil {
		log.Printf("Error beginning database transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(db.Ctx)

	_, err = tx.Exec(db.Ctx, "DELETE FROM user_roles WHERE user_id = $1;", deleteUser.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(db.Ctx, "DELETE FROM users WHERE id = $1;", deleteUser.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(db.Ctx)
	if err != nil {
		log.Printf("Error commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
