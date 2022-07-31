package main

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gorilla/mux"
	"linum-banking-api/nordigen"
	"net/http"
)

const (
	AuthTokenKey = "auth_token"
)

func initializeWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/test", testHandler)
	r.HandleFunc("/nordigen/institutions", nordigenGetInstitutionsHandler)
	r.HandleFunc("/nordigen/requisition-link", nordigenRequisitionLinkHandler)
	r.Handle("/nordigen/requisition-link", authorize(http.HandlerFunc(nordigenRequisitionLinkHandler)))
	http.Handle("/", r)
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := verifyToken(idToken, getFirebaseInstance())
		ctx := context.WithValue(r.Context(), AuthTokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("Authorization")
	if idToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	verifyToken(idToken, getFirebaseInstance())
}

func nordigenGetInstitutionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := nordigen.GetToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	// TODO: Pass Country Code to request
	institutions, err := nordigen.GetInstitutionsForCountry("de", token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	jsonResp, err := json.Marshal(institutions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func nordigenRequisitionLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, ok := r.Context().Value(AuthTokenKey).(*auth.Token)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Remove
	fmt.Println(authToken.UID)

	var reqRequest createRequisitionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// TODO: Perhaps even check userId at this point

	token, err := nordigen.GetToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	requisition, err := nordigen.CreateRequisition(reqRequest.InstitutionId, reqRequest.UserLanguage, reqRequest.RedirectUrl, token)
	if err != nil {
		// TODO: Differentiate between errors
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	// TODO: Store reference
	err = storeRequisitionId(requisition.Id, authToken.UID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	jsonResp, err := json.Marshal(createRequisitionResponse{Link: requisition.Link})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
