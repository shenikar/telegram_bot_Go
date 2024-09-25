package server

import (
	"encoding/json"
	"net/http"
	"telegram_bot_go/service"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.GetLast24HoursRequests()
	if err != nil {
		http.Error(w, "Error retrieving stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
