package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/controllers"
)

func AdminRoutes() *http.ServeMux {

	router := http.NewServeMux()

	// user
	router.HandleFunc("GET /get-all-users", controllers.GetAllUsers)
	router.HandleFunc("GET /get-single-user/{userID}", controllers.GetSingleUser)

	return router

}
