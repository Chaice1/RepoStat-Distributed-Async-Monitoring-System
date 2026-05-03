package collectordomain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	CreatedAt   string
}

type Subscription struct {
	Id        uuid.UUID
	RepoName  string
	OwnerName string
	CreatedAt string
}

var (
	ErrNotFound      = errors.New("NOT_FOUND")
	ErrInternalError = errors.New("INTERNAL_ERROR")
	ErrBadRequest    = errors.New("BAD_REQUEST")
)

type GitHubClient interface {
	GetRepoInfo(context.Context, string, string) (*RepoInfo, error)
}
