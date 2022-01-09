package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cemtanrikut/go-api-horsea/api"
	"github.com/cemtanrikut/go-api-horsea/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	FirstName  string    `json:"firstname", bson:"firstname"`
	LastName   string    `json:"lastname", bson:"lastname" `
	Email      string    `json:"email", bson:"email"`
	Password   string    `json:"password", bson:"password"`
	CreateDate time.Time `json:"createdate", bson:"createdate"`
	UpdateDate time.Time `json:"updatedate", bson:"updatedate"`
	IsDeleted  bool      `json:"isdeleted", bson:"isdeleted"`
}

func SignUp(resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(req.Body).Decode(&user)
	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))
	user.CreateDate = time.Now()
	user.IsDeleted = false

	checkEmail := CheckEmail(user.Email, client, collection)
	if checkEmail {
		return helper.ReturnResponse(http.StatusUnauthorized, "", "This mail address is already exist.")

	}

	_, insertErr := collection.InsertOne(context.Background(), user)
	if insertErr != nil {
		return helper.ReturnResponse(http.StatusBadRequest, "", insertErr.Error())
	}

	jsonResult, jsonError := json.Marshal(user)
	if jsonError != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", jsonError.Error())

	}

	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")

}

func LogIn(resp http.ResponseWriter, req *http.Request, client *mongo.Client, ctx context.Context, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user, dbUser User

	json.NewDecoder(req.Body).Decode(&user)

	err := collection.FindOne(context.Background(), bson.M{"email": user.Email, "isdeleted": false}).Decode(&dbUser)

	if err != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())

	}

	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))

	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	fmt.Println(userPass, dbPass)
	fmt.Println(user.Password, dbUser.Password)

	//passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	res := bytes.Equal(userPass, dbPass)
	if !res {
		return helper.ReturnResponse(http.StatusBadRequest, "", err.Error())
	}

	jsonResult, jsonError := json.Marshal(user.Email)
	if jsonError != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
	}
	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")

}

func GetUser(email string, resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user User

	userData := collection.FindOne(context.Background(), bson.M{"email": email, "isdeleted": false})
	err := userData.Decode(&user)

	if err != nil {
		return helper.ReturnResponse(http.StatusNotFound, "", err.Error())
	}

	jsonResult, jsonError := json.Marshal(user)
	if jsonError != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
	}

	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")
}

func GetUsers(client *mongo.Client, resp http.ResponseWriter, req *http.Request, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var userMList []primitive.M

	cursor, err := collection.Find(context.Background(), bson.M{"isdeleted": false})
	if err != nil {
		return helper.ReturnResponse(http.StatusNotFound, "fdsdfgtgsfgssfg", err.Error())
	}

	for cursor.Next(context.Background()) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
		}
		userMList = append(userMList, user)
	}
	defer cursor.Close(context.Background())

	jsonResult, err := json.Marshal(userMList)
	if err != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
	}

	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")

}

func UpdateUser(resp http.ResponseWriter, req *http.Request, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user User

	json.NewDecoder(req.Body).Decode(&user)

	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))

	updatedData, updateErr := collection.UpdateOne(context.Background(), bson.M{"email": user.Email, "isdeleted": false}, bson.D{{"$set",
		bson.D{
			{"firstname", user.FirstName},
			{"lastname", user.LastName},
			{"password", user.Password},
			{"updatedate", time.Now()},
		},
	}})
	if updateErr != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", updateErr.Error())
	}
	jsonResult, err := json.Marshal(updatedData)
	if err != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
	}

	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")

}

func UpdatePassword(resp http.ResponseWriter, req *http.Request, collection *mongo.Collection, email string) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user User

	userData := collection.FindOne(context.Background(), bson.M{"email": email, "isdeleted": false})
	err := userData.Decode(&user)

	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))

	if err != nil {
		return helper.ReturnResponse(http.StatusNotFound, "", err.Error())
	}

	updatedData, updateErr := collection.UpdateOne(context.Background(), bson.M{"email": user.Email, "isdeleted": false}, bson.D{{"$set",
		bson.D{
			{"password", user.Password},
		},
	}})

	if updateErr != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", updateErr.Error())
	}
	jsonResult, err := json.Marshal(updatedData)
	if err != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "", err.Error())
	}

	return helper.ReturnResponse(http.StatusOK, string(jsonResult), "")
}

func DeleteUser(email string, resp http.ResponseWriter, req *http.Request, collection *mongo.Collection) api.Response {
	resp.Header().Set("Content-Type", "application/json")
	var user User

	json.NewDecoder(req.Body).Decode(&user)

	_, err := collection.UpdateOne(context.Background(), bson.M{"email": email, "isdeleted": false}, bson.D{{"$set",
		bson.D{
			{"isdeleted", true},
		},
	}})
	if err != nil {
		return helper.ReturnResponse(http.StatusInternalServerError, "Something went wrong to Deleting User :(", err.Error())
	}

	return helper.ReturnResponse(http.StatusOK, "User deleted successfully!", "")
}

func CheckEmail(email string, client *mongo.Client, collection *mongo.Collection) bool {
	var dbUser User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&dbUser)
	fmt.Println("data - ", err)
	if err == nil {
		fmt.Println(email, " already exist")
		return true
	}
	return false
}
