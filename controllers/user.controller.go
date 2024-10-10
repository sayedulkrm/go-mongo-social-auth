package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/sayedulkrm/go-mongo-social-auth/config"
	"github.com/sayedulkrm/go-mongo-social-auth/helpers"
	"github.com/sayedulkrm/go-mongo-social-auth/models"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var validate = validator.New()

func UserRegister(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Entering UserRegister")

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to parse multipart form data")
		return
	}

	// Get form values
	email := r.FormValue("email")
	password := r.FormValue("password")
	userName := r.FormValue("userName")

	if email == "" || password == "" || userName == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Please provide all required fields")
		return
	}

	newUserData := helpers.RegisterUserDataStruct{
		Email:    email,
		Password: password,
		UserName: userName,
	}

	validationErr := validate.Struct(newUserData)
	if validationErr != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to validate user data")

		return
	}

	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancle()

	var existingUser models.USER

	err = config.UserCollection.FindOne(ctx, bson.M{"email": newUserData.Email}).Decode(&existingUser)

	if err == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "User Email Already Exists")
		return
	}

	err = config.UserCollection.FindOne(ctx, bson.M{"user_name": newUserData.UserName}).Decode(&existingUser)

	if err == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "User Name Already Exists")
		return
	}

	// Generating TOken
	token, activationCode, err := helpers.CreateActivationToken(newUserData)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to create activation token")
		return
	}

	mailData := fmt.Sprintf("Name: %s, Activation Code: %s", newUserData.UserName, activationCode)

	fmt.Println("Sending activation email with data:", mailData) // This would be your mail sending logic

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Email sent to %s . Please activate your account", newUserData.Email),
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Activate User

// Struct for the data within the JWT token
type ActivationPayload struct {
	UserData       helpers.RegisterUserDataStruct `json:"userData"`
	ActivationCode string                         `json:"activationCode"`
}

func verifyActivationToken(tokenString string) (*ActivationPayload, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACTIVATION_SECRECT")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		// Extract data from claims
		payload := &ActivationPayload{
			UserData: helpers.RegisterUserDataStruct{
				UserName: (*claims)["userData"].(map[string]interface{})["userName"].(string),
				Email:    (*claims)["userData"].(map[string]interface{})["email"].(string),
				Password: (*claims)["userData"].(map[string]interface{})["password"].(string),
			},
			ActivationCode: (*claims)["activationCode"].(string),
		}
		return payload, nil
	}

	return nil, fmt.Errorf("invalid token claims")

}

func ActivateUser(w http.ResponseWriter, r *http.Request) {

	logrus.Info("Entering Activate User")

	var reqBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to parse body === Activate User File")

		return
	}

	activationCode := reqBody["activationCode"]
	activationToken := reqBody["activationToken"]

	// Verify the activation token

	activationPayload, err := verifyActivationToken(activationToken)
	logrus.Info(activationPayload)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Invalid activation token === Err Comes")

		return
	}

	// Compare the activation codes
	if activationPayload.ActivationCode != activationCode {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Invalid activation token === Token not mathc")

		return
	}

	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancle()

	var existingUser models.USER

	err = config.UserCollection.FindOne(ctx, bson.M{"email": activationPayload.UserData.Email}).Decode(&existingUser)

	if err == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "User Email Already Exists")

		return
	}

	err = config.UserCollection.FindOne(ctx, bson.M{"user_name": activationPayload.UserData.UserName}).Decode(&existingUser)

	if err == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "User Name Already Exists")

		return
	}

	validationErr := validate.Struct(activationPayload.UserData)
	if validationErr != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to validate user data")

		return
	}

	hashedPassword, err := helpers.HashPassword(activationPayload.UserData.Password)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to hash password")
		return
	}

	newUser := models.USER{
		UserName: activationPayload.UserData.UserName,
		Email:    activationPayload.UserData.Email,
		Password: hashedPassword,
	}

	_, err = config.UserCollection.InsertOne(ctx, newUser)

	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed To create User")

		return
	}

	// Send success response
	response := map[string]interface{}{
		"success": true,
		"message": "Account activated successfully. Please login",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

type UserLoginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var userLoginData UserLoginData

	err := json.NewDecoder(r.Body).Decode(&userLoginData)

	if err != nil {
		utils.LogError(r, err)
		utils.ErrorResponse(w, r, http.StatusBadRequest, "failed to decode")
		return
	}

	validationErr := validate.Struct(userLoginData)
	if validationErr != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed to validate user data")

		return
	}
	// Find the user by email
	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancle()

	// Find the user by email
	var existingUser models.USER

	err = config.UserCollection.FindOne(ctx, bson.M{"email": userLoginData.Email}).Decode(&existingUser)

	if err != nil {
		utils.LogError(r, err)
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid email")
		return
	}

	isPasswordMatch := helpers.ComparePassword(existingUser.Password, userLoginData.Password)
	// Check if the password is correct
	if !isPasswordMatch {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, "Invalid password")
		return
	}

	utils.SendToken(existingUser, 200, "Welcome", w, r)

}

// logout
func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User logout logic goes here"))
}

// =================== 		Profile  	==============

// Get user profile
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get user profile logic goes here"))
}

// Update user profile http.ResponseWriter, r *http.Request

// =================== 		Social Auth  	==============

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

	ctx, cancle := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancle()
	// Find the user by email
	var existingUser models.USER

	err = config.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

	// 	If the user is found, the FindOne function will decode the result into existingUser, and err will be nil.
	// If no user is found, err will be set to mongo.ErrNoDocuments, and you should allow the code to proceed, as this means the username does not already exist.

	if err == nil {
		utils.SendToken(existingUser, 200, "Welcome", w, r)
		return
	}

	// If the user does not exist, create a new username by combining first and last names
	firstName := user.FirstName
	lastName := user.LastName

	uniqueUsername, err := helpers.EnsureUniqueUsername(ctx, firstName, lastName)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, "Failed to ensure unique username")
		return
	}

	// Create a new user with the unique username and other details
	newUser := models.USER{
		UserName:  uniqueUsername,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		UserImage: models.UserImage{
			URL: user.AvatarURL, // Save the Avatar URL
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}
	// Insert the new user into the database
	_, err = config.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "Failed To create User")
		return
	}

	// Send token to the new user
	utils.SendToken(newUser, 201, fmt.Sprintf("Account created successfully. Welcome, %s!", newUser.FirstName), w, r)

}

// ========================
