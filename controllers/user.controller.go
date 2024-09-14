package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/sayedulkrm/go-mongo-social-auth/config"
)

var userCollection = config.OpenCollection(config.CreatedMongoClient, "user")

func UserRegister(w http.ResponseWriter, r *http.Request) {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("user register")
	// }

	w.Write([]byte("User registration logic goes here"))
}

// func UserLogin()

type key int

const ProviderKey key = 0

// Middleware to set provider into context
func SetProviderInContext(r *http.Request, provider string) *http.Request {
	ctx := context.WithValue(r.Context(), ProviderKey, provider)
	return r.WithContext(ctx)
}

// Extract provider from context
func GetProviderFromContext(r *http.Request) string {
	provider, ok := r.Context().Value(ProviderKey).(string)
	if !ok {
		return ""
	}
	return provider
}

// Social Auth

func HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	// // Extract the provider (e.g., "google")

	provider := r.PathValue("provider") // Extracting 'google' from '/auth/google'

	// Add provider as a query parameter
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	// Begin the authentication process
	gothic.BeginAuthHandler(w, r)
}

func GetGoogleAuthCallbackFunc(w http.ResponseWriter, r *http.Request) {

	// Complete the authentication process
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// Log the authenticated user
	fmt.Println(user.Name)
	fmt.Println(user.Email)
	fmt.Println(user.AvatarURL)
	fmt.Println(user.FirstName)
	fmt.Println(user.LastName)
	fmt.Println(user.NickName)

	// Redirect the user after authentication
	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
}
