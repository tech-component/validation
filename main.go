package main

import (
	"log"
	"net/http"

	"github.com/mauleyzaola/validation/middlewares"
	"github.com/mauleyzaola/validation/rest"
	"github.com/mauleyzaola/validation/validators"
)

func main() {
	http.HandleFunc("/users", createUserHandler())

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createUserHandler() http.HandlerFunc {
	handler := rest.NewServer(validators.NewValidator())
	return middlewares.MethodChecker(
		http.MethodPost,
		middlewares.JSONValidator(
			handler.Validator(),
			handler.CreateUser,
		),
	)
}
