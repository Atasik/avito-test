package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"segmenter/internal/domain"
)

type addUserInput struct {
	ID               int              `json:"id"`
	SegmentsToDelete []domain.Segment `json:"segmentsToDelete"`
	SegmentsToAdd    []domain.Segment `json:"segmentsToAdd"`
}

// @Summary Add user to segment
// @Tags user
// @ID	 add-user-to-segment
// @Accept json
// @Product json
// @Param   input body addUserInput true "user id, segments to delete/add"
// @Success	200		    {string}	string     "updated"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/segment/user [post]
func (h *Handler) AddUserToSegment(w http.ResponseWriter, r *http.Request) {
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

	var inp addUserInput
	err = json.Unmarshal(body, &inp)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}

	err = h.Services.UpsertUser(inp.ID, inp.SegmentsToAdd, inp.SegmentsToDelete)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refactor, remove hardcode
	resp, err := json.Marshal(map[string]interface{}{
		"updated": "done",
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

// @Summary Get User segments
// @Tags user
// @ID	 get-user-segments
// @Accept json
// @Product json
// @Param   input body domain.User true "user id"
// @Success	200		    {object}	getSegmentsResponse     "segments"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/user [post]
func (h *Handler) GetUserSegments(w http.ResponseWriter, r *http.Request) {
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

	var user domain.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}
	segments, err := h.Services.GetUserSegments(user.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetSegmentsResponse(w, segments, http.StatusOK)
}
