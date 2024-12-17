package rest

import (
	"encoding/json"
	"net/http"

	"github.com/mauleyzaola/validation/domain"
	"github.com/mauleyzaola/validation/interfaces"
)

type Server struct {
	validator interfaces.Validator
}

func NewServer(validator interfaces.Validator) *Server {
	return &Server{
		validator: validator,
	}
}

func (s *Server) Validator() interfaces.Validator {
	return s.validator
}

func (s *Server) CreateUser(user domain.User, w http.ResponseWriter, _ *http.Request) {
	// TODO: save the user to a database
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
