package subscriber_domain

import (
	"github.com/google/uuid"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	CreatedAt   string
}

type RepoInfoTaskMessage struct {
	Id        uuid.UUID
	Payload   RepoInfoTaskMessagePayload
	CreatedAt string
}

type RepoInfoTaskMessagePayload struct {
	Repo  string
	Owner string
}
