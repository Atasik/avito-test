package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"segmenter/internal/domain"
)

type addUserInput struct {
	ID               int              `json:"id" example:"1" validate:"required,gt=0"`
	SegmentsToDelete []domain.Segment `json:"segmentsToDelete"`
	SegmentsToAdd    []domain.Segment `json:"segmentsToAdd"`
}

// @Summary Add user to segment
// @Tags user
// @ID	 add-user-to-segment
// @Accept json
// @Product json
// @Param   input body addUserInput true "user id, segments to delete/add"
// @Success	200		    {object}	statusResponse
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/segment/user [post]
func (h *Handler) addUserToSegment(w http.ResponseWriter, r *http.Request) {
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
		log.Println("kekw", err.Error())
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(inp)
	if err != nil {
		newErrorResponse(w, "bad input", http.StatusBadRequest)
		return
	}

	err = h.Services.User.UpsertUserSegments(inp.ID, inp.SegmentsToAdd, inp.SegmentsToDelete)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(statusResponse{"done"})
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
func (h *Handler) getUserSegments(w http.ResponseWriter, r *http.Request) {
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

	err = h.Validator.Struct(user)
	if err != nil {
		newErrorResponse(w, "bad input", http.StatusBadRequest)
		return
	}

	segments, err := h.Services.User.GetSegments(user.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetSegmentsResponse(w, segments, http.StatusOK)
}
