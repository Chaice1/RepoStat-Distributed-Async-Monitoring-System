package processor_controller_errors

import (
	"errors"
	"log/slog"
	processor_domain "repo-stat/processor/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleErrorFromDomainToGRPC(err error, log *slog.Logger, action string) error {
	switch {
	case errors.Is(err, processor_domain.ErrBadRequest):
		log.Error("invalid argument", "action", action, "error", err)
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	case errors.Is(err, processor_domain.ErrInternalError):
		log.Error("internal error", "action", action, "error", err)
		return status.Error(codes.Internal, codes.Internal.String())
	case errors.Is(err, processor_domain.ErrNotFound):
		log.Error("not found", "action", action, "error", err)
		return status.Error(codes.NotFound, codes.NotFound.String())
	default:
		log.Error("unknown error", "action", action, "error", err)
		return status.Error(codes.Unknown, codes.Unknown.String())
	}
}
