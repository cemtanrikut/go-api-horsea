package main

import (
	"context"
	"log"
	"net/http"

	api "github.com/cemtanrikut/go-api-horsea/api/user"
	db "github.com/cemtanrikut/go-api-horsea/db"
	"github.com/cemtanrikut/go-api-horsea/helper"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
)

var client *mongo.Client
var ctx context.Context
var userCollection *mongo.Collection
var router *mux.Router

func main() {
	log.Println("Starting the application")

	router, ctx, client, userCollection = db.MongoClient("user_collection")

	router.HandleFunc("/api/user/login", LogIn).Methods(http.MethodPost)
	router.HandleFunc("/api/user/signup", SignUp).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{email}", GetUser).Methods(http.MethodGet)
	router.HandleFunc("/api/user/users", GetUsers).Methods(http.MethodPost)
	router.HandleFunc("/api/user/update", UpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/delete/{email}", DeleteUser).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	result := api.LogIn(w, r, client, ctx, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	result := api.SignUp(w, r, client, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	result := api.GetUser(email, w, r, client, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	result := api.GetUsers(client, w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	result := api.UpdateUser(w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	email := v["email"]
	result := api.DeleteUser(email, w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
