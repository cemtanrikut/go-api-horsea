package router

import (
	"context"
	"log"
	"net/http"

	api "github.com/cemtanrikut/go-api-horsea/api/user"
	"github.com/cemtanrikut/go-api-horsea/db"
	"github.com/cemtanrikut/go-api-horsea/helper"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var userCollection *mongo.Collection
var router *mux.Router

func MuxUserHandler() {
	router, ctx, client, userCollection = db.MongoClient("user_collection")

	router.HandleFunc("/api/user/login", logIn).Methods(http.MethodPost)
	router.HandleFunc("/api/user/signup", signUp).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{email}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/api/user/users", getUsers).Methods(http.MethodPost)
	router.HandleFunc("/api/user/update", updateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user/changepassword/{email}", changePassword).Methods(http.MethodPost)
	router.HandleFunc("/api/user/delete/{email}", deleteUser).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	result := api.SignUp(w, r, client, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func logIn(w http.ResponseWriter, r *http.Request) {
	result := api.LogIn(w, r, client, ctx, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	result := api.GetUser(email, w, r, client, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func getUsers(w http.ResponseWriter, r *http.Request) {
	result := api.GetUsers(client, w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func updateUser(w http.ResponseWriter, r *http.Request) {
	result := api.UpdateUser(w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func changePassword(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	result := api.ChangePassword(w, r, userCollection, email)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	email := v["email"]
	result := api.DeleteUser(email, w, r, userCollection)
	byteRes := helper.JsonMarshal(result)
	w.Write(byteRes)
}
