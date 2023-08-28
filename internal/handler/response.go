package handler

import (
	"encoding/json"
	"net/http"
	"segmenter/internal/domain"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getSegmentsResponse struct {
	Data []domain.Segment `json:"segments"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(errorResponse{msg}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newGetSegmentsResponse(w http.ResponseWriter, segments []domain.Segment, status int) {
	resp, _ := json.Marshal(getSegmentsResponse{segments}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}
