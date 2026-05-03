package processor_domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	Status      string
	CreatedAt   time.Time
}

type Ping struct {
	Reply string
}

type GetRepoInfoRequestBody struct {
	Repo  string
	Owner string
}

var (
	ErrNotFound      = errors.New("NOT_FOUND")
	ErrInternalError = errors.New("INTERNAL_ERROR")
	ErrBadRequest    = errors.New("BAD_REQUEST")
)

type CollectorClient interface {
	Ping(context.Context) (*Ping, error)
}

type RepoInfoTaskMessage struct {
	Id        uuid.UUID
	Payload   RepoInfoTaskMessagePayload
	CreatedAt time.Time
}

type RepoInfoTaskMessagePayload struct {
	Repo  string
	Owner string
}

type RepoInfoResMessage struct {
	Id        uuid.UUID
	Payload   []byte
	CreatedAt time.Time
}

type Repository interface {
	GetRepoInfo(context.Context, string) (*RepoInfo, error)
	CreateFetchingTaskTransaction(context.Context, string, string) error
}
