package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth/gothic"
	"github.com/sayedulkrm/go-mongo-social-auth/config"
	"github.com/sayedulkrm/go-mongo-social-auth/helpers"
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

func GetGoogleAuthCallbackFunc(w http.ResponseWriter, r *http.Request) {

	helpers.SocialAuthHelper()

	// We have to check if we are geting {provider} from params

	value := r.PathValue("provider")
	fmt.Println("value", value)

	// Extract the provider from the URL path
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	provider := pathSegments[2] // Extracting 'google' from '/auth/google/callback'

	fmt.Println("provider", provider)

	// // Set provider in request context
	// r = SetProviderInContext(r, provider)

	// // Check if the provider is correctly set
	// fmt.Println("provider", GetProviderFromContext(r))

	// get the user from the session

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println(user)

	// redirect the user
	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)

}

// func HandleProviderLogin(w http.ResponseWriter, r *http.Request) {

// 	// get the provider from the request

// 	provider := r.PathValue("provider")

// 	fmt.Println("provider", provider)

// 	// redirect the user to the provider's login page

// 	gothic.BeginAuthHandler(w, r)

// }
