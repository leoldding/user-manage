package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoldding/user-manage/database"
)

func (db *DB) getAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Admin getting all users")

	// get database connection from pool
	conn, err := db.Pool.Acquire(db.Ctx)
	if err != nil {
		log.Printf("Error acquiring connection from pool: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(db.Ctx, "SELECT username, first_name, last_name FROM users;")
	if err != nil {
		log.Printf("Error retrieving users from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []*database.User
	for rows.Next() {
		var user database.User
		err = rows.Scan(&user.Username, &user.FirstName, &user.Password)
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

func (db *DB) updateUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Admin updating user with id %s", id)

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
	_, err = conn.Exec(db.Ctx, "UPDATE users SET username = $2, password = $3, first_name = $4, last_name = %5 WHERE id= $1;", id, updateUser.Username, updateUser.Password, updateUser.FirstName, updateUser.LastName)
	if err != nil {
		log.Printf("Error updating user values: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *DB) deleteUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Admin deleting user with id %s", id)

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
	_, err = tx.Exec(db.Ctx, "DELET FROM user_roles WHERE user_id = $1;", id)
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
