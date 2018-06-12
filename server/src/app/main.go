package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"context"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func public(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello public!\n"))
}

func private(w http.ResponseWriter, r *http.Request) {
	opt := option.WithCredentialsFile(os.Getenv("CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	auth, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	print(idToken)
	token, err := auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
	w.Write([]byte("hello private!\n"))
}

func main() {
	port := os.Getenv("SERVER_PORT")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	r := mux.NewRouter()
	r.HandleFunc("/public", public)
	r.HandleFunc("/private", private)

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods)(r)))
}
