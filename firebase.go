package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"log"
	"sync"
)

var lock = &sync.Mutex{}

var firebaseApp *firebase.App

func getFirebaseInstance() *firebase.App {
	if firebaseApp == nil {
		lock.Lock()
		defer lock.Lock()
		if firebaseApp == nil {
			var err error
			firebaseApp, err = initializeFirebase()
			if err != nil {
				log.Fatalf("Could not setup firebase connection")
			}
		}
	}
	return firebaseApp
}

func initializeFirebase() (*firebase.App, error) {
	// opt := option.WithCredentialsFile("./linum-dev-firebase-adminsdk-b7u2y-1faae2b9fe.json")
	// app, err := firebase.NewApp(context.Background(), nil, opt)
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}

func verifyToken(idToken string, app *firebase.App) *auth.Token {
	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	return token
}
