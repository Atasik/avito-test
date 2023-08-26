package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"segmenter/internal/repository"
)

type addUserInput struct {
	Id               int                  `json:"id"`
	SegmentsToDelete []repository.Segment `json:"segmentsToDelete"`
	SegmentsToAdd    []repository.Segment `json:"segmentsToAdd"`
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
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	err = h.Services.UpsertUser(inp.Id, inp.SegmentsToAdd, inp.SegmentsToDelete)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
// @Param   input body repository.User true "user id"
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

	var user repository.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}
	segments, err := h.Services.GetUserSegments(user.Id)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetSegmentsResponse(w, segments, http.StatusOK)
}
