package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/controllers"
	"github.com/sayedulkrm/go-mongo-social-auth/helpers"
)

func AdminRoutes() *http.ServeMux {

	router := http.NewServeMux()

	// user
	router.HandleFunc("GET /get-all-users", (helpers.AuthorizeRoles("admin"))(controllers.GetAllUsers))
	router.HandleFunc("GET /get-single-user/{userID}", (helpers.AuthorizeRoles("admin"))(controllers.GetSingleUser))

	return router

}
