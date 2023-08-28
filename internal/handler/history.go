package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"segmenter/pkg/timejson"
	"time"
)

type getReportInput struct {
	ID     int                    `json:"id"`
	Period timejson.YearMonthTime `json:"period"`
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
		newErrorResponse(w, "can't unpack payload", http.StatusBadRequest)
		return
	}

	reportName, err := h.Services.CreateReport(time.Time(inp.Period), inp.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refactor, remove hardcode, create cfg
	resp, err := json.Marshal(map[string]interface{}{
		"link": r.Host + "/reports/" + reportName,
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
