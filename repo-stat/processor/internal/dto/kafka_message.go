package processor_dto

import (
	"github.com/google/uuid"
)

type RepoInfoTaskMessage struct {
	Id        uuid.UUID                  `json:"id"`
	Payload   RepoInfoTaskMessagePayload `json:"payload"`
	CreatedAt string                     `json:"created_at"`
}

type RepoInfoTaskMessagePayload struct {
	Repo  string `json:"repo"`
	Owner string `json:"owner"`
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
