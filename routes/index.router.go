package routes

import (
	"net/http"

	"github.com/sayedulkrm/go-mongo-social-auth/middlewares"
)


func SetupRoutes() http.Handler {


	mux := http.NewServeMux()

	mux.HandleFunc("GET /",middlewares.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<h1>Server is working. To See Frontend <a href="http://localhost:3000"> Click Here </a></h1>`
		w.Write([]byte(html))
	}))

	return mux

}