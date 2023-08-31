package handler

import (
	"net/http"
	"segmenter/internal/service"

	_ "segmenter/docs"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

const (
	appJSON = "application/json"
)

type Handler struct {
	Services  *service.Service
	Validator *validator.Validate
}

func (h *Handler) InitRoutes(reportsDir string) http.Handler {
	r := mux.NewRouter()

	// TODO: remove hardcode
	staticHandler := http.StripPrefix(
		"/reports/",
		http.FileServer(http.Dir("./reports")),
	)

	r.PathPrefix("/reports/").Handler(staticHandler)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/api/segment", h.createSegment).Methods("POST")
	r.HandleFunc("/api/segment", h.deleteSegment).Methods("DELETE")
	r.HandleFunc("/api/segment/user", h.addUserToSegment).Methods("POST")
	r.HandleFunc("/api/user", h.getUserSegments).Methods("POST")
	r.HandleFunc("/api/user/history", h.getReport).Methods("POST")

	mux := panicMiddleware(r)

	return mux
}
