package helpers

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

var validate = validator.New()

// func VerifyPassword()

// func HashPassword()

func SocialAuthHelper() {
	err := godotenv.Load()
	if err != nil {
		utils.NewErrorHandler("Failed To load env", http.StatusBadGateway)
		panic(err)

	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// Have to use store- dont know why

	sessionSecret := os.Getenv("SESSION_SECRET")

	key := sessionSecret // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	// Callback URL has to be exact same in Google cloud console

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8000/api/v1/auth/google/callback"),
	)
}
