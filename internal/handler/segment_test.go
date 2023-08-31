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

func TestHandler_createSegment(t *testing.T) {
	type mockBehaviour func(r *service_mocks.MockSegment, input domain.Segment)

	tests := []struct {
		name                 string
		inputBody            string
		input                domain.Segment
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "TEST"}`,
			input: domain.Segment{
				Name: "TEST",
			},
			mockBehaviour: func(r *service_mocks.MockSegment, input domain.Segment) {
				r.EXPECT().Create(input).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Bad Input",
			inputBody: `{"name": ""}`,
			input: domain.Segment{
				Name: "",
			},
			mockBehaviour:        func(r *service_mocks.MockSegment, input domain.Segment) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad input"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"name":"TEST"}`,
			input: domain.Segment{
				Name: "TEST",
			},
			mockBehaviour: func(r *service_mocks.MockSegment, input domain.Segment) {
				r.EXPECT().Create(input).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoSegment := service_mocks.NewMockSegment(c)
			test.mockBehaviour(repoSegment, test.input)

			services := &service.Service{Segment: repoSegment}

			validate := validator.New()
			h := &Handler{
				Services:  services,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/segment", h.createSegment).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/segment",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteSegment(t *testing.T) {
	type mockBehaviour func(r *service_mocks.MockSegment, input domain.Segment)

	tests := []struct {
		name                 string
		inputBody            string
		input                domain.Segment
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"TEST"}`,
			input: domain.Segment{
				Name: "TEST",
			},
			mockBehaviour: func(r *service_mocks.MockSegment, input domain.Segment) {
				r.EXPECT().Delete(input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"done"}`,
		},
		{
			name:      "Bad Input",
			inputBody: `{"name": ""}`,
			input: domain.Segment{
				Name: "",
			},
			mockBehaviour:        func(r *service_mocks.MockSegment, input domain.Segment) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"bad input"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"name":"TEST"}`,
			input: domain.Segment{
				Name: "TEST",
			},
			mockBehaviour: func(r *service_mocks.MockSegment, input domain.Segment) {
				r.EXPECT().Delete(input).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoSegment := service_mocks.NewMockSegment(c)
			test.mockBehaviour(repoSegment, test.input)

			services := &service.Service{Segment: repoSegment}

			validate := validator.New()
			h := &Handler{
				Services:  services,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/segment", h.deleteSegment).Methods("DELETE")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/segment",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
