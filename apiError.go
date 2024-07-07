package swift

import (
	"net/http"

	res "github.com/brownhounds/swift/response"
)

func ApiError(w http.ResponseWriter, s int) {
	res.Json(w, s, res.Map{
		"status":  s,
		"message": http.StatusText(s),
	})
}
