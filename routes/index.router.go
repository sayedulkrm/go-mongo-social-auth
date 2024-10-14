package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

func SetupRoutes() http.Handler {

	rootRoutes := http.NewServeMux()
	// Providing socail auth helper

	// Set up user routes
	userRoutes := UserRoutes()
	// admin
	adminRoutes := AdminRoutes()

	// Set up admin routes
	rootRoutes.Handle("/api/v1/user/", middlewares.LogMiddleware(http.StripPrefix("/api/v1/user", userRoutes).(http.HandlerFunc)))
	rootRoutes.Handle("/api/v1/admin/", middlewares.LogMiddleware(http.StripPrefix("/api/v1/admin", adminRoutes).(http.HandlerFunc)))

	rootRoutes.HandleFunc("/", middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic(utils.NewErrorHandler("Path Not Found", http.StatusBadRequest))
	}))

	rootRoutes.HandleFunc("/{$}", middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<h1>Server is working. To See Frontend <a href="http://localhost:3000"> Click Here </a></h1>`
		w.Write([]byte(html))
	}))

	return rootRoutes

}
