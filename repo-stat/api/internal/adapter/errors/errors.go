package adapter_errors

import (
	"log/slog"
	"repo-stat/api/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandleFromGRPCToDomainWithLog(err error, log *slog.Logger, action string) error {
	status, ok := status.FromError(err)
	if !ok {
		log.Error("internal error", "action", "get GRPC status of error", "error", domain.ErrInternalError)
		return domain.ErrInternalError
	}

	switch status.Code() {
	case codes.InvalidArgument:
		log.Error("invalid argument", "action", action, "error", domain.ErrBadRequest)
		return domain.ErrBadRequest
	case codes.NotFound:
		log.Error("not found", "action", action, "error", domain.ErrNotFound)
		return domain.ErrNotFound
	case codes.AlreadyExists:
		log.Error("subscription already exists", "action", action, "error", domain.ErrSubscriptionAlreadyExists)
		return domain.ErrSubscriptionAlreadyExists
	default:
		log.Error("internal error", "action", action, "error", domain.ErrInternalError)
		return domain.ErrInternalError
	}
}
