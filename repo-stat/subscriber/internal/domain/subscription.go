package subscriber_domain

import "github.com/google/uuid"

type Subscription struct {
	Id        uuid.UUID
	RepoName  string
	OwnerName string
	CreatedAt string
}
