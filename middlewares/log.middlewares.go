package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)


func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Request Received: %s %s", r.Method, r.URL.Path)

		next(w,r)
	}
}