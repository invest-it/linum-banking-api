package main

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"linum-banking-api/src/nordigen"
	"log"
	"net/http"
	"time"
)

const (
	AuthTokenKey = "auth_token"
)

func initializeWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/test", testHandler)
	r.HandleFunc("/nordigen/institutions", nordigenGetInstitutionsHandler)
	// r.HandleFunc("/nordigen/requisition-link", nordigenRequisitionLinkHandler)
	r.HandleFunc("/nordigen/callback", nordigenCallbackHandler)
	r.Handle("/nordigen/requisition-link", authorize(http.HandlerFunc(nordigenRequisitionLinkHandler)))
	r.Handle("/nordigen/transactions/{requisitionId}", authorize(http.HandlerFunc(nordigenLoadTransactionsHandler)))
	r.Handle("/nordigen/transactions", authorize(http.HandlerFunc(nordigenLoadAllTransactionsHandler)))
	http.Handle("/", r)
	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("No IdToken was found in the request header")
			return
		}
		token := verifyToken(idToken, getFirebaseInstance())
		ctx := context.WithValue(r.Context(), AuthTokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	message := "Hello world"
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, message)
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
		fmt.Print("Error while fetching token: ")
		fmt.Println(err)
		return
	}
	// TODO: Pass Country Code to request
	institutions, err := nordigen.GetInstitutionsForCountry("de", token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("Error while fetching institutions: ")
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
		fmt.Println("AuthToken not found!")
		return
	}

	var reqRequest CreateRequisitionRequest
	if err := json.NewDecoder(r.Body).Decode(&reqRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	token, err := nordigen.GetToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("Error while fetching token: ")
		fmt.Println(err)
		return
	}

	requisition, err := nordigen.CreateRequisition(reqRequest.InstitutionId, reqRequest.UserLanguage, reqRequest.RedirectUrl, token)
	if err != nil {
		// TODO: Differentiate between errors
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error while creating requisition: ")
		nordigen.HandleApiError(err)
		return
	}
	fmt.Println(requisition.Reference)
	err = storeRequisitionId(requisition.Id, requisition.Reference, authToken.UID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	jsonResp, err := json.Marshal(CreateRequisitionResponse{Link: requisition.Link})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// FLOW: Get Accounts and their Transactions from RequisitionId
// -> Parse them into custom models
// -> Save them in Firebase
// -> Store date of first and last entry in db

func nordigenLoadTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	authToken, ok := r.Context().Value(AuthTokenKey).(*auth.Token)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	urlVars := mux.Vars(r)
	reqId, ok := urlVars["requisitionId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the id belongs to the user
	if !userHasRequisition(reqId, authToken.UID) {
		w.WriteHeader(http.StatusUnauthorized) // TODO: Handle case where no item was found
		fmt.Println("User does not have requisition")
		return
	}

	// TODO: Perhaps better handle token internally
	token, err := nordigen.GetToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	requisition, err := nordigen.GetRequisitionById(reqId, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, accountId := range requisition.Accounts {
		transactions, err := nordigen.GetTransactionsForAccountId(accountId, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return or break? EXAMPLE: json->error: Could not fetch for AccountId {ID}
			return
		}
		jsonData, err := json.Marshal(transactions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return or break? EXAMPLE: json->error: Could not fetch for AccountId {ID}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		break
		// -> PARSE into custom model
		// -> Save in firebase
		// -> Store last transaction date
	}

}

func nordigenLoadAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {

}

func getLatestTransaction(transactions []nordigen.Transaction) *nordigen.Transaction {
	if len(transactions) == 0 {
		return nil
	}
	latest := transactions[0]
	for _, transaction := range transactions {
		transactionDate, err := time.Parse("2006-01-02", transaction.ValueDate)
		latestDate, err := time.Parse("2006-01-02", latest.ValueDate)
		if err != nil {
			// TODO: Handle error -> Throw or quiet?
			return &latest
		}
		if transactionDate.Unix() > latestDate.Unix() {
			latest = transaction
		}
	}
	return &latest
}

func nordigenCallbackHandler(w http.ResponseWriter, r *http.Request) {
	refIds, ok := r.URL.Query()["ref"]
	if !ok || len(refIds[0]) < 1 {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Could not finish the request, a required reference is missing.")
	}

	refId := refIds[0]

	err := updateRequisitionState(refId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Successfully approved requisition")
}
