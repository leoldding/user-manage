package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/leoldding/user-manage/database"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) login(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging In")

	// decode request body into user
	var creds *database.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
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

	// retrieve id of user
	var id []byte
	err = conn.QueryRow(db.Ctx, "SELECT id FROM users WHERE username = $1;", creds.Username).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("User does not exist: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			log.Printf("Error getting id from database: %v", err)
		}
	}

	// retrieve stored password from database
	var storedPass []byte
	err = conn.QueryRow(db.Ctx, "SELECT password FROM users WHERE id = $1;", id).Scan(&storedPass)
	if err != nil {
		log.Printf("Error getting stored password from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// compare stored password hash and submitted password
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(creds.Password))
	if err != nil {
		log.Printf("Incorrect password for user: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// retrieve role of user
	var role string
	err = conn.QueryRow(db.Ctx, "SELECT name FROM roles WHERE id = (SELECT role_id FROM user_roles WHERE user_id = $1);", id).Scan(&role)
	if err != nil {
		log.Printf("Error getting user's role from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create jwt claims
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing JWT: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set jwt int http only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user-jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Write([]byte(role))
}

func (db *DB) logout(w http.ResponseWriter, r *http.Request) {
	log.Println("User Logging Out")

	// invalidate jwt token
	http.SetCookie(w, &http.Cookie{
		Name:     "user-jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}

func isAuthenticated(w http.ResponseWriter, r *http.Request) {
	// retrieve jwt
	cookie, err := r.Cookie("user-jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("Missing JWT: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			log.Printf("Invalid JWT: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}

	// parse token
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Printf("Incorrect JWT signing method: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return nil, errors.New("Incorrect signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// invalid token
	if !token.Valid {
		log.Printf("Invalid JWT: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
