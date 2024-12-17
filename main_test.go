package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mauleyzaola/validation/mocks"
	"github.com/mauleyzaola/validation/validators"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mauleyzaola/validation/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// assume mocked repository always succeeds
	repository := &mocks.RepositoryMock{
		CreateUserFunc: func(ctx context.Context, user domain.User) (string, bool, error) {
			return "55624c8b-c0b6-4562-a51b-68a61c73e3c6", true, nil
		},
	}
	validator := validators.NewValidator()
	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "valid user",
			method: http.MethodPost,
			body: domain.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"55624c8b-c0b6-4562-a51b-68a61c73e3c6","email":"test@example.com","password":"password123"}`,
		},
		{
			name:   "invalid email",
			method: http.MethodPost,
			body: domain.User{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field Email: invalid email format",
		},
		{
			name:   "password too short",
			method: http.MethodPost,
			body: domain.User{
				Email:    "test@example.com",
				Password: "short",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field Password: minimum length is 8",
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			body:           domain.User{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				bodyBytes []byte
				err       error
			)

			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.body)
				require.NoError(t, err, "failed to marshal request body")
			}

			req := httptest.NewRequest(tt.method, "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := createUserHandler(repository, validator)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")

			responseBody := strings.TrimSpace(rr.Body.String())
			assert.Equal(t, tt.expectedBody, responseBody, "handler returned unexpected body")
		})
	}
}
