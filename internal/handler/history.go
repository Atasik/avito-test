package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type getReportInput struct {
	Id      int       `json:"id"`
	BeginAt time.Time `json:"begin_at"`
	EndAt   time.Time `json:"end_at"`
}

// @Summary Get Report
// @Tags history
// @ID	 get-report
// @Accept json
// @Product json
// @Param   input body getReportInput true "user id, interval begin and end"
// @Success	200		    {string}	string     "link"
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/user/history [post]
func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
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

	var inp getReportInput
	err = json.Unmarshal(body, &inp)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	reportName, err := h.Services.CreateReport(inp.BeginAt, inp.EndAt, inp.Id)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"link": "http://localhost:8080/reports/" + reportName,
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
