package helper

import (
	"net/http"

	"github.com/cemtanrikut/go-api-horsea/api"
)

func ReturnResponse(statusCode int, user string, err string) api.Response {
	statusText := http.StatusText(statusCode)
	return api.Response{
		Data:         user,
		StatusCode:   statusCode,
		ErrorMessage: statusText + " - " + err,
	}
}
