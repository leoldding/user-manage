package main

import (
	"context"
	"log"
	"net/http"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/leoldding/user-manage/database"
	"github.com/leoldding/user-manage/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	ctx := context.Background()
	pool, err := database.NewDatabase(ctx)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}
	defer pool.Close()

	handlers.RegisterHandlers(router)

	headersOk := gHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := gHandlers.AllowedOrigins([]string{"localhost:3000"})
	methodsOk := gHandlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PATCH"})
	credentialsOk := gHandlers.AllowCredentials()

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", gHandlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(router)))
}
