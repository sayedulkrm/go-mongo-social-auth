package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/helpers"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}

func GetSingleUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userID")

	user, err := helpers.GetUserDetailsById(userId)

	if err != nil {
		utils.NewErrorHandler("Failed to find user", http.StatusNotFound)
	}

	// Set response header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set status code to 200

	response := map[string]interface{}{
		"success": true,
		"message": "User found",
		"user":    user,
	}

	// Encode the response to JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
