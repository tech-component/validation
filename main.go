package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/mauleyzaola/validation/domain"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/mauleyzaola/validation/assets"
	"github.com/mauleyzaola/validation/middlewares"
	"github.com/mauleyzaola/validation/rest"
	"github.com/mauleyzaola/validation/validators"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// connect to db
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()
	db, err := NewPostgresDB(ctx, dbUrl)
	if err != nil {
		return err
	}
	defer func() { db.Close() }()

	// apply db migrations
	if err = migrateDb(assets.MigrationFiles(),
		"files/migrations", dbUrl); err != nil {
		return err
	}

	// TODO: move this code somewhere else
	id, ok, err := createUser(ctx, db.pool, domain.User{
		Email:    "mauricio.leyzaola@gmail.com",
		Password: "super-secret",
	})
	if err != nil {
		return err
	}
	log.Printf("generated user id: %s ok: %v\n", id, ok)

	// define handlers
	http.HandleFunc("/users", createUserHandler())

	// start http server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}

type PostgresDB struct {
	pool *pgxpool.Pool
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

func migrateDb(files fs.FS, filePath, dbURL string) error {
	driver, err := iofs.New(files, filePath)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

func connectPG(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewPostgresDB(ctx context.Context, dbUrl string) (*PostgresDB, error) {
	pool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}

func createUser(ctx context.Context, pool *pgxpool.Pool, user domain.User) (string, bool, error) {
	var (
		id string
		ok bool
	)

	query := `SELECT * FROM create_user($1, $2);`

	// Execute the query
	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(&id, &ok)
	if err != nil {
		return "", false, fmt.Errorf("failed to create user: %w", err)
	}

	return id, ok, nil
}
