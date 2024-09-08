package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
	"github.com/sayedulkrm/go-mongo-social-auth/utils"
)

func UserRoutes() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html")
		html := `<h1>Register. To See Frontend <a href="http://localhost:3000"> Click Here </a></h1>`
		w.Write([]byte(html))

	})

	router.HandleFunc("/", middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic(utils.NewErrorHandler("Path Not Found", http.StatusBadRequest))
	}))

	return router

}
