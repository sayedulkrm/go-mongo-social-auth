package controllers

import "net/http"

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}

func GetSingleUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get single user"))
}
