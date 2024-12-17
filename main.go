package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tech-component/validation/assets"
	"github.com/tech-component/validation/database"
	"github.com/tech-component/validation/interfaces"
	"github.com/tech-component/validation/middlewares"
	"github.com/tech-component/validation/migrations"
	"github.com/tech-component/validation/repositories"
	"github.com/tech-component/validation/rest"
	"github.com/tech-component/validation/validators"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("PORT")
	// connect to db
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()

	// apply db migrations
	if err := migrations.MigrateDb(assets.MigrationFiles(),
		"files/migrations", dbUrl); err != nil {
		return err
	}

	// instantiate repository
	db, err := database.NewPostgresDB(ctx, dbUrl)
	if err != nil {
		return err
	}
	defer func() { db.Close() }()
	repository := repositories.NewPGRepository(db.Pool())
	validator := validators.NewValidator()

	// define handlers
	http.HandleFunc("/users", createUserHandler(repository, validator))

	// start http server
	log.Println("Server starting on :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	return nil
}

func createUserHandler(repository interfaces.Repository, validator interfaces.Validator) http.HandlerFunc {
	// TODO: move this function somewhere else, put here for the sake of simplicity
	handler := rest.NewServer(repository, validator)
	return middlewares.MethodChecker(
		http.MethodPost,
		middlewares.JSONValidator(
			handler.Validator(),
			handler.CreateUser,
		),
	)
}
