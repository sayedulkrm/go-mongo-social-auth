package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/controllers"
	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

func UserRoutes() *http.ServeMux {

	router := http.NewServeMux()
	// Auth
	router.HandleFunc("POST /register", controllers.UserRegister)
	router.HandleFunc("POST /activate", controllers.ActivateUser)
	router.HandleFunc("POST /login", controllers.UserLogin)

	// Social Auth
	router.HandleFunc("GET /auth/{provider}", controllers.HandleProviderLogin)
	router.HandleFunc("GET /auth/{provider}/callback", controllers.GetGoogleAuthCallbackFunc)
	// From frontend We need to call it continuosly. To get user from sesson
	router.HandleFunc("GET /google/login/success", controllers.CreateUserFromSocalAuth)

	// Authienticated Routes
	router.HandleFunc("GET /me", controllers.GetUserProfile)

	router.HandleFunc("/", middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic(utils.NewErrorHandler("Path Not Found", http.StatusBadRequest))
	}))

	return router

}
