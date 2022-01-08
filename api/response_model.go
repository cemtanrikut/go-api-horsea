package api

type Response struct {
	Data         string `json:"data"`
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}
