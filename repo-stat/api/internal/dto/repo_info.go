package dto

type RepoInfo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stargazers  uint64 `json:"stars"`
	Forks       uint64 `json:"forks"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type GetRepositoriesInfoResponse struct {
	Repositories []*RepoInfo `json:"repositories"`
}
