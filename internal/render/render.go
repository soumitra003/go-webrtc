package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type baseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RenderError(w http.ResponseWriter, err error) {
	baseResponse := baseResponse{
		Status:  "500",
		Message: "ERROR",
		Data:    err.Error(),
	}
	RenderJSON(w, baseResponse)
}

func RenderBaseResponse(w http.ResponseWriter, data interface{}) {
	baseResponse := baseResponse{
		Status:  "200",
		Message: "SUCCESS",
		Data:    data,
	}

	RenderJSON(w, baseResponse)
}

func RenderJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func RenderStreamMessage(w http.ResponseWriter, msg []byte) {
	fmt.Fprintf(w, "data: Message: %s\n\n", msg)
}

func Redirect(w http.ResponseWriter, url string, status int) {
	w.Header().Set("Location", url)
	w.WriteHeader(status)
}
