package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayedulkrm/go-mongo-social-auth/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	ID          primitive.ObjectID `json:"_id"`
	Email       string             `json:"email"`
	FirstName   string             `json:"first_name"`
	UserName    string             `json:"user_name"`
	LastName    string             `json:"last_name"`
	PhoneNumber string             `json:"phone_number"`
	UserRole    string             `json:"user_role"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	UserImage   models.UserImage   `json:"user_image"`
}

// Set access and refresh token expiry
var accessTokenExpire = 300   // in minutes
var refreshTokenExpire = 1200 // in minutes

// Set JWT secret from environment
var accessTokenSecret = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
var refreshTokenSecret = os.Getenv("JWT_REFRESH_TOKEN_SECRET")

func SignAccessToken() (string, error) {

	var user models.USER

	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * time.Duration(accessTokenExpire)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(accessTokenSecret))
}

func SignRefreshToken() (string, error) {

	var user models.USER
	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Minute * time.Duration(refreshTokenExpire)).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(refreshTokenSecret))

}

func SendToken(user models.USER, statusCode int, message string, w http.ResponseWriter, r *http.Request) {

	accessToken, err := SignAccessToken()
	if err != nil {
		logrus.Error("Failed To Generate Access Token")
		ErrorResponse(w, r, http.StatusBadRequest, "Failed To generate Access Token")
		return
	}

	refreshToken, err := SignRefreshToken()
	if err != nil {
		logrus.Error("Failed To Generate Refresh Token")
		ErrorResponse(w, r, http.StatusBadRequest, "Failed To generate Refresh Token")
		return
	}

	userResponse := UserResponse{
		ID:          user.Id,
		Email:       user.Email,
		FirstName:   user.FirstName,
		UserName:    user.UserName,
		LastName:    user.LastName,
		PhoneNumber: user.Phone_Number,
		UserRole:    user.UserRole,
		CreatedAt:   user.Created_At,
		UpdatedAt:   user.Updated_At,
		UserImage:   user.UserImage,
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(accessTokenExpire)),
		MaxAge:   accessTokenExpire * 60,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenExpire)),
		MaxAge:   refreshTokenExpire * 60 * 60 * 24,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"user":    userResponse,
		"success": true,
	})

}

// Forsical auth
func SendTokenAndRedirect(user models.USER, w http.ResponseWriter, r *http.Request) {

	accessToken, err := SignAccessToken()
	if err != nil {
		logrus.Error("Failed To Generate Access Token")
		http.Error(w, "Failed to generate access token", http.StatusBadRequest)
		return
	}

	refreshToken, err := SignRefreshToken()
	if err != nil {
		logrus.Error("Failed To Generate Refresh Token")
		http.Error(w, "Failed to generate refresh token", http.StatusBadRequest)
		return
	}

	// Create cookies for the tokens
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(accessTokenExpire)),
		MaxAge:   accessTokenExpire * 60,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenExpire)),
		MaxAge:   refreshTokenExpire * 60 * 60 * 24,
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteNoneMode,
	}

	// Set the cookies
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	// Redirect to the main page (frontend)
	http.Redirect(w, r, "http://localhost:3000/", http.StatusSeeOther)
}
