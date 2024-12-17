package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/tech-component/validation/interfaces"
)

func JSONValidator[T any](validator interfaces.Validator, next func(T, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload T
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		if err := validator.ValidateStruct(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next(payload, w, r)
	}
}
