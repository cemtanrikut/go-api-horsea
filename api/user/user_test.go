package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cemtanrikut/go-api-horsea/db"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	var user User
	user.FirstName = "Cem"
	user.LastName = "Tanrikut"
	user.Email = "cem@gmail.com"
	user.Password = "1234"

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(user)

	r, _ := http.NewRequest("POST", "/api/user/signup", io.Reader(reqBodyBytes))
	w := httptest.NewRecorder()

	w.Body.Write(reqBodyBytes.Bytes())

	_, _, client, collection := db.MongoClient("user_collection")

	response := SignUp(w, r, client, collection)

	assert.Equal(t, http.StatusOK, response.StatusCode)

}

func TestLogIn(t *testing.T) {

}

func TestGetUser(t *testing.T) {

}

func TestGetUsers(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}

func TestDeleteUser(t *testing.T) {

}

func TestCheckEmail(t *testing.T) {

}
