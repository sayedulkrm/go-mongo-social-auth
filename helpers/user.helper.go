package helpers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sayedulkrm/go-mongo-social-auth/config"
	"github.com/sayedulkrm/go-mongo-social-auth/models"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// 4 digit
func GenerateActivationToken() string {

	RandomCrypto, _ := rand.Prime(rand.Reader, 12)

	return fmt.Sprintf("%d", RandomCrypto)

}

// Generate a username by combining first and last names with a 4-digit number
func GenerateUsername(firstName, lastName string) (string, error) {
	// Ensure that firstName and lastName are not empty
	if firstName == "" && lastName == "" {
		return "", fmt.Errorf("first name and last name are both empty")
	}

	// Create base username by combining first and last names
	baseUsername := firstName + lastName

	// Generate a random 4-digit number and append it to the base username
	randomNumber := GenerateActivationToken()
	logrus.Info("GENERATED USER NAME: ", baseUsername+randomNumber)

	return baseUsername + randomNumber, nil
}

// Ensure the username is unique
func EnsureUniqueUsername(ctx context.Context, firstName, lastName string) (string, error) {
	var existingUser models.USER

	// Keep generating usernames until a unique one is found
	for {
		// Generate a new username with a random number
		baseUsername, err := GenerateUsername(firstName, lastName)
		if err != nil {
			return "", err
		}

		err = config.UserCollection.FindOne(ctx, bson.M{"user_name": baseUsername}).Decode(&existingUser)
		if err != nil {
			// If err is mongo.ErrNoDocuments, the username is unique
			if err == mongo.ErrNoDocuments {
				return baseUsername, nil
			}
			return "", err
		}
		// If a user with the same username is found, loop and regenerate
	}
}

// type User struct {
//     FirstName string `json:"firstName"`
//     LastName  string `json:"lastName"`
//     Email     string `json:"email"`
//     Password  string `json:"password"`
//     File      string `json:"file"` // File upload handling could vary
// }

type RegisterUserDataStruct struct {
	UserName string `json:"userName" validate:"required,min=3,max=20"` // Ensure username is between 3-20 characters
	Email    string `json:"email" validate:"required,email"`           // Ensure a valid email is provided
	Password string `json:"password" validate:"required,min=6"`
}

// Create Activation Token
func CreateActivationToken(user RegisterUserDataStruct) (string, string, error) {
	activationCode := GenerateActivationToken()

	// Create token claims using HS256
	claims := jwt.MapClaims{
		"userData":       user,
		"activationCode": activationCode,
		"exp":            time.Now().Add(time.Minute * 15).Unix(), // 15 min expiration
	}

	logrus.Info("JWT Activation Secret: ", os.Getenv("JWT_ACTIVATION_SECRECT"))

	// Use HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACTIVATION_SECRECT")))
	if err != nil {
		logrus.Errorf("Token creation failed: %v", err) // Log the actual error
		return "", "", fmt.Errorf("failed to create token: %v", err)
	}

	return tokenString, activationCode, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("Failed to hash password")
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GetUserDetailsById retrieves a user by their ID
func GetUserDetailsById(userID string) (*models.USER, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.USER

	// Find the user by ID and decode into the user variable
	err := config.UserCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, utils.NewErrorHandler("Failed to find user", http.StatusNotFound)
	}

	return &user, nil
}

// AuthorizeRoles middleware to check user roles from request body
func AuthorizeRoles(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Parse the request body to extract user roles
			var userRoles string
			if err := json.NewDecoder(r.Body).Decode(&userRoles); err != nil {
				utils.NewErrorHandler("Failed to decode user roles", http.StatusBadRequest)
				return
			}

			for _, role := range allowedRoles {
				if role == userRoles {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Role does not match, return forbidden
			utils.NewErrorHandler("Forbidden - Insufficient privileges", http.StatusForbidden)
		}
	}
}

// Social auth

func SocialAuthHelper() {
	err := godotenv.Load()
	if err != nil {
		utils.NewErrorHandler("Failed To load env", http.StatusBadGateway)
		panic(err)

	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	key := os.Getenv("SESSION_SECRET")
	if key == "" {
		logrus.Fatal("SESSION_SECRET is not set in the environment")
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(86400 * 30)      // 30 days
	store.Options.HttpOnly = true // HttpOnly should be enabled for security
	store.Options.Secure = true   // Ensure this is false in development (set to true in production)
	// Required if you are dealing with cross-domain issues
	gothic.Store = store

	// Callback URL has to be exact same in Google cloud console

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8000/api/v1/user/auth/google/callback", "profile", "email"),
	)
}
