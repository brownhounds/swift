package response

import (
	"encoding/json"
	"net/http"
)

type Map map[string]interface{}

func Json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err.Error())
	}
}
