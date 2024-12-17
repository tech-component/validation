package rest

import (
	"encoding/json"
	"net/http"

	"github.com/mauleyzaola/validation/domain"
	"github.com/mauleyzaola/validation/interfaces"
)

type Server struct {
	repository interfaces.Repository
	validator  interfaces.Validator
}

func NewServer(repository interfaces.Repository, validator interfaces.Validator) *Server {
	return &Server{
		repository: repository,
		validator:  validator,
	}
}

func (s *Server) Validator() interfaces.Validator {
	return s.validator
}

func (s *Server) CreateUser(user domain.User, w http.ResponseWriter, r *http.Request) {
	// save the user to a database
	id, _, err := s.repository.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
