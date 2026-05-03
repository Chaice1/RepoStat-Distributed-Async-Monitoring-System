package domain

import "github.com/google/uuid"

type Subscription struct {
	Id        uuid.UUID
	OwnerName string
	RepoName  string
	CreatedAt string
}
