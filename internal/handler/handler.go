package handler

import (
	"net/http"
	"segmenter/internal/service"

	_ "segmenter/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

const (
	appJSON = "application/json"
)

type Handler struct {
	Services *service.Service
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	// TODO: remove hardcode
	staticHandler := http.StripPrefix(
		"/reports/",
		http.FileServer(http.Dir("./reports")),
	)

	r.PathPrefix("/reports/").Handler(staticHandler)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/api/segment", h.CreateSegment).Methods("POST")
	r.HandleFunc("/api/segment", h.DeleteSegment).Methods("DELETE")
	r.HandleFunc("/api/segment/user", h.AddUserToSegment).Methods("POST")
	r.HandleFunc("/api/user", h.GetUserSegments).Methods("POST")
	r.HandleFunc("/api/user/history", h.GetReport).Methods("POST")

	mux := panicMiddleware(r)

	return mux
}
