package api_controller_errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
)

func HandleErrorsFromDomainToHTTP(w http.ResponseWriter, err error, log *slog.Logger, action string) {
	var code int
	switch {
	case errors.Is(err, domain.ErrBadRequest):
		log.Error("invalid argument", "action", action, "error", err)
		code = http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		log.Error("not found", "action", action, "error", err)
		code = http.StatusNotFound
	case errors.Is(err, domain.ErrSubscriptionAlreadyExists):
		log.Error("subscription already exists", "action", action, "error", err)
		code = http.StatusConflict
	default:
		log.Error("internal error", "action", action, "error", err)
		code = http.StatusInternalServerError
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		log.Error("failed to write errorInfo", "error", err)
	}

}
