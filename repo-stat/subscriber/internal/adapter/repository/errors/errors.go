package subscriber_repo_errors

import (
	subscriber_domain "repo-stat/subscriber/internal/domain"

	"github.com/jackc/pgx/v5"
)

func HandleErrorFromDBToDomain(err error) error {
	switch err {
	case pgx.ErrNoRows:
		return subscriber_domain.ErrNotFound
	default:
		return subscriber_domain.ErrInternalError
	}
}
