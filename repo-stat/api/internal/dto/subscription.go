package dto

import "github.com/google/uuid"

type Subscription struct {
	Id        uuid.UUID `json:"id"`
	RepoName  string    `json:"repo_name"`
	OwnerName string    `json:"owner_name"`
	CreatedAt string    `json:"created_at"`
}

type CreateSubscriptionRequest struct {
	RepoName  string `json:"repo_name"`
	OwnerName string `json:"owner_name"`
}

type GetSubscriptionsResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}
