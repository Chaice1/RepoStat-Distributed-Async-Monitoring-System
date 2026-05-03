package grpc_errors

import (
	subscriber_domain "repo-stat/subscriber/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleErrorFromDomainToGRPC(err error) error {

	switch err {
	case subscriber_domain.ErrBadRequest:
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	case subscriber_domain.ErrNotFound:
		return status.Error(codes.NotFound, codes.NotFound.String())
	case subscriber_domain.ErrSubscriptionAlreadyExists:
		return status.Error(codes.AlreadyExists, codes.AlreadyExists.String())
	case subscriber_domain.ErrInternalError:
		return status.Error(codes.Internal, codes.Internal.String())
	default:
		return status.Error(codes.Unknown, codes.Unknown.String())
	}
}
