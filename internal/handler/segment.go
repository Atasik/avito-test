package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"segmenter/internal/domain"
)

// @Summary Create segment
// @Tags segment
// @ID	 create-segment
// @Accept json
// @Product json
// @Param   input body domain.Segment true "Segment content"
// @Success	200		    {integer}	integer     "id"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/segment [post]
func (h *Handler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var segment domain.Segment
	err = json.Unmarshal(body, &segment)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}
	id, err := h.Services.CreateSegment(segment)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refactor, remove hardcode
	resp, err := json.Marshal(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}

// @Summary Delete segment
// @Tags segment
// @ID	 delete-segment
// @Accept json
// @Product json
// @Param   input body domain.Segment true "Segment content"
// @Success	200		    {integer}	integer     "id"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/segment [delete]
func (h *Handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var segment domain.Segment
	err = json.Unmarshal(body, &segment)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}
	err = h.Services.DeleteSegment(segment)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refactor, remove hardcode
	resp, err := json.Marshal(map[string]interface{}{
		"deleted": "done",
	})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}
