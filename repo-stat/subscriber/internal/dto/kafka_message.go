package subscriber_dto

import "github.com/google/uuid"

type DeleteSubscriptionMessage struct {
	Id          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stargazers  uint64    `json:"stargazers_count"`
	Forks       uint64    `json:"forks"`
	CreatedAt   string    `json:"created_at"`
	Error       string    `json:"error,omitempty"`
}
