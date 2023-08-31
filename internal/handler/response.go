package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"segmenter/internal/domain"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getSegmentsResponse struct {
	Data []domain.Segment `json:"segments"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, err := json.Marshal(errorResponse{msg}) //nolint:errcheck
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newGetSegmentsResponse(w http.ResponseWriter, segments []domain.Segment, status int) {
	resp, err := json.Marshal(getSegmentsResponse{segments}) //nolint:errcheck
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}
