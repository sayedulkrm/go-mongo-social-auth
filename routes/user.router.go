package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/controllers"
	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

func UserRoutes() *http.ServeMux {

	router := http.NewServeMux()

	// router.HandleFunc("POST /register", controllers.UserRegister)

	// Social Auth
	router.HandleFunc("GET /auth/{provider}/callback", controllers.GetGoogleAuthCallbackFunc)
	// router.HandleFunc("GET /auth/{provider}", controllers.HandleProviderLogin)

	router.HandleFunc("/", middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic(utils.NewErrorHandler("Path Not Found", http.StatusBadRequest))
	}))

	return router

}
