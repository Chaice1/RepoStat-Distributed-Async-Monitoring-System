package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// NewPingHandler
// @Summary      Проверка статуса системы
// @Description  Возвращает статус доступности подписчка и процессора
// @Tags         system
// @Produce      json
// @Success      200  {object}  dto.PingResponse
// @Failure      503  {object}  dto.PingResponse
// @Router       /api/ping [get]
func NewPingHandler(log *slog.Logger, ping *usecase.Ping, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collector_status := ping.Execute(r.Context())
		processor_status := agu.Ping(r.Context())
		w.Header().Set("Content-Type", "application/json")
		var response *dto.PingResponse
		if collector_status == domain.PingStatusUp && processor_status == domain.PingStatusUp {
			response = CreatePingResponce("ok", processor_status, collector_status)
			w.WriteHeader(http.StatusOK)
		} else {
			response = CreatePingResponce("degraded", processor_status, collector_status)
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}

	}
}

func CreatePingResponce(system_status string, processor_status domain.PingStatus, subscriber_status domain.PingStatus) *dto.PingResponse {

	return &dto.PingResponse{
		Status: system_status,
		Services: []dto.ServicesInfo{
			dto.ServicesInfo{
				Name:   "processor",
				Status: string(processor_status),
			},
			dto.ServicesInfo{
				Name:   "subscriber",
				Status: string(subscriber_status),
			},
		},
	}
}
