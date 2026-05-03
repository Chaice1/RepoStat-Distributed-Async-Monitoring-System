package collectorrespmodel

import (
	"github.com/google/uuid"
)

type RepoInfo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stargazers  uint64 `json:"stargazers_count"`
	Forks       uint64 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type RepoInfoResMessage struct {
	Id          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stargazers  uint64    `json:"stargazers_count"`
	Forks       uint64    `json:"forks"`
	CreatedAt   string    `json:"created_at"`
	Error       string    `json:"error,omitempty"`
}

type RepoInfoMessage struct {
	Id        uuid.UUID              `json:"id"`
	Payload   RepoInfoMessagePayload `json:"payload"`
	CreatedAt string                 `json:"created_at"`
}

type RepoInfoMessagePayload struct {
	Repo  string `json:"repo"`
	Owner string `json:"owner"`
}
