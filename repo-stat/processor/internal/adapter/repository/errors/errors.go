package processor_db_errors

import (
	"errors"
	"log/slog"
	processor_domain "repo-stat/processor/internal/domain"

	"github.com/jackc/pgx/v5"
)

func ErrorHandleFromDBToDomain(err error, log *slog.Logger, action string) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Error("no rows found", "action", action, "error", err)
		return processor_domain.ErrNotFound
	default:
		log.Error("internal error", "action", action, "error", err)
		return processor_domain.ErrInternalError
	}
}
