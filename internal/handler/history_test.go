package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"segmenter/internal/service"
	service_mocks "segmenter/internal/service/mocks"
	"segmenter/pkg/timejson"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_getReports(t *testing.T) {
	type mockBehaviour func(r *service_mocks.MockHistory, input getReportInput)

	tests := []struct {
		name                 string
		inputBody            string
		input                getReportInput
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id":1, "period":"2023-08"}`,
			input: getReportInput{
				ID: 1,
				Period: timejson.YearMonthTime{
					Time: time.Date(
						2023, 8, 1, 00, 00, 00, 0, time.UTC),
				},
			},
			mockBehaviour: func(r *service_mocks.MockHistory, input getReportInput) {
				r.EXPECT().CreateReport(input.Period.Time, input.ID).Return("link", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"link":"example.com/reports/link"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"id":1, "period":"2023-08"}`,
			input: getReportInput{
				ID: 1,
				Period: timejson.YearMonthTime{
					Time: time.Date(
						2023, 8, 1, 00, 00, 00, 0, time.UTC),
				},
			},
			mockBehaviour: func(r *service_mocks.MockHistory, input getReportInput) {
				r.EXPECT().CreateReport(input.Period.Time, input.ID).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		{
			name:      "Bad Payload",
			inputBody: `{"id":1, "period":"wrong input"}`,
			input: getReportInput{
				ID: 1,
			},
			mockBehaviour:        func(r *service_mocks.MockHistory, input getReportInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"can't unpack payload"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoHistory := service_mocks.NewMockHistory(c)
			test.mockBehaviour(repoHistory, test.input)

			services := &service.Service{History: repoHistory}

			validate := validator.New()
			h := &Handler{
				Services:  services,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/user/history", h.getReport).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/history",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
