package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	api "github.com/cemtanrikut/go-api-horsea/api/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

var client *mongo.Client
var ctx context.Context
var collection *mongo.Collection

func main() {
	log.Println("Starting the application")

	router := mux.NewRouter()
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:LCtfPjhpm1am7HRd@sandbox.0sac2.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	collection = client.Database("horsea_db").Collection("user_collection")

	router.HandleFunc("/api/user/login", LogIn).Methods(http.MethodPost)
	router.HandleFunc("/api/user/signup", SignUp).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{email}", GetUser).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	result := api.LogIn(w, r, client, ctx, collection)
	jsonResult, jsonError := json.Marshal(result)
	if jsonError != nil {
		w.Write([]byte(jsonError.Error()))
	}
	w.Write(jsonResult)
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	result := api.SignUp(w, r, client, collection)
	jsonResult, jsonError := json.Marshal(result)
	if jsonError != nil {
		w.Write([]byte(jsonError.Error()))
	}
	w.Write(jsonResult)
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	result := api.GetUser(email, w, r, client, collection)
	jsonResult, jsonError := json.Marshal(result)
	if jsonError != nil {
		w.Write([]byte(jsonError.Error()))
	}
	w.Write(jsonResult)
}
