package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/tech-component/validation/interfaces"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tech-component/validation/domain"
	"github.com/tech-component/validation/mocks"
	"github.com/tech-component/validation/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		body       interface{}
		method     string
		repository interfaces.Repository
	}
	validator := validators.NewValidator()
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid user",
			args: args{
				body: domain.User{
					Email:    "test@example.com",
					Password: "password123",
				},
				method: http.MethodPost,
				repository: &mocks.RepositoryMock{
					CreateUserFunc: func(ctx context.Context, user domain.User) (string, bool, error) {
						return "55624c8b-c0b6-4562-a51b-68a61c73e3c6", true, nil
					},
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"55624c8b-c0b6-4562-a51b-68a61c73e3c6","email":"test@example.com","password":"password123"}`,
		},
		{
			name: "invalid email",
			args: args{
				method: http.MethodPost,
				body: domain.User{
					Email:    "invalid-email",
					Password: "password123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field Email: invalid email format",
		},
		{
			name: "password too short",
			args: args{
				method: http.MethodPost,
				body: domain.User{
					Email:    "test@example.com",
					Password: "short",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field Password: minimum length is 8",
		},
		{
			name: "invalid method",
			args: args{
				method: http.MethodGet,
				body:   domain.User{},
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name: "invalid json",
			args: args{
				method: http.MethodPost,
				body:   "invalid json",
			},
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

			switch v := tt.args.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.args.body)
				require.NoError(t, err, "failed to marshal request body")
			}

			req := httptest.NewRequest(tt.args.method, "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := createUserHandler(tt.args.repository, validator)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")

			responseBody := strings.TrimSpace(rr.Body.String())
			assert.Equal(t, tt.expectedBody, responseBody, "handler returned unexpected body")
		})
	}
}
