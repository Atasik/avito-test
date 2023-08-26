package handler

import (
	"encoding/json"
	"net/http"
	"segmenter/internal/repository"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getSegmentsResponse struct {
	Data []repository.Segment `json:"segments"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(errorResponse{msg})
	w.WriteHeader(status)
	w.Write(resp)
}

func newGetSegmentsResponse(w http.ResponseWriter, segments []repository.Segment, status int) {
	resp, _ := json.Marshal(getSegmentsResponse{segments})
	w.WriteHeader(status)
	w.Write(resp)
}
