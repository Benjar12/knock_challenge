package httpresp

import (
	"encoding/json"
	"net/http"
)

type ReturnValue struct {
	Status int
	Error  string
	Data   interface{}
}

// WriteJSON makes the responses sent to the client concestent.
// StatusCode defaults to 200 unless an error happens.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp, _ := json.Marshal(ReturnValue{
			Status: statusCode,
			Error:  err.Error(),
		})

		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
		return
	}

	resp, _ := json.Marshal(ReturnValue{
		Status: statusCode,
		Data:   data,
	})

	w.WriteHeader(statusCode)
	w.Write([]byte(resp))
}
