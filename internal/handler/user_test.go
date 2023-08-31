package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"segmenter/internal/domain"
	"segmenter/internal/service"
	service_mocks "segmenter/internal/service/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_addUserToSegment(t *testing.T) {
	type mockBehaviour func(r *service_mocks.MockUser, input addUserInput)

	tests := []struct {
		name                 string
		inputBody            string
		input                addUserInput
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1,"segmentsToAdd":[{"name":"TEST"}]}`,
			input: addUserInput{
				ID: 1,
				SegmentsToAdd: []domain.Segment{
					{
						Name: "TEST",
					},
				},
			},
			mockBehaviour: func(r *service_mocks.MockUser, input addUserInput) {
				r.EXPECT().UpsertUserSegments(input.ID, input.SegmentsToAdd, input.SegmentsToDelete).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"done"}`,
		},
		{
			name:      "Bad Input",
			inputBody: `{"id":-1}`,
			input: addUserInput{
				ID: -1,
			},
			mockBehaviour:        func(r *service_mocks.MockUser, input addUserInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad input"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"id":1,"segmentsToAdd":[{"name":"TEST"}]}`,
			input: addUserInput{
				ID: 1,
				SegmentsToAdd: []domain.Segment{
					{
						Name: "TEST",
					},
				},
			},
			mockBehaviour: func(r *service_mocks.MockUser, input addUserInput) {
				r.EXPECT().UpsertUserSegments(input.ID, input.SegmentsToAdd, input.SegmentsToDelete).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoUser := service_mocks.NewMockUser(c)
			test.mockBehaviour(repoUser, test.input)

			services := &service.Service{User: repoUser}

			validate := validator.New()
			h := &Handler{
				Services:  services,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/segment/user", h.addUserToSegment).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/segment/user",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getUserSegments(t *testing.T) {
	type mockBehaviour func(r *service_mocks.MockUser, input domain.User)

	tests := []struct {
		name                 string
		inputBody            string
		input                domain.User
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1}`,
			input:     domain.User{ID: 1},
			mockBehaviour: func(r *service_mocks.MockUser, input domain.User) {
				r.EXPECT().GetSegments(input.ID).Return([]domain.Segment{{Name: "TEST"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"segments":[{"name":"TEST"}]}`,
		},
		{
			name:                 "Bad Input",
			inputBody:            `{"id":-1}`,
			input:                domain.User{ID: -1},
			mockBehaviour:        func(r *service_mocks.MockUser, input domain.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"id":1}`,
			input:     domain.User{ID: 1},
			mockBehaviour: func(r *service_mocks.MockUser, input domain.User) {
				r.EXPECT().GetSegments(input.ID).Return([]domain.Segment{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoUser := service_mocks.NewMockUser(c)
			test.mockBehaviour(repoUser, test.input)

			services := &service.Service{User: repoUser}

			validate := validator.New()
			h := &Handler{
				Services:  services,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/user", h.getUserSegments).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
