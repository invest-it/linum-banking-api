package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func initializeWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/test", TestHandler)
	http.Handle("/", r)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("Authorization")
	if idToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	verifyToken(idToken, getFirebaseInstance())
}
