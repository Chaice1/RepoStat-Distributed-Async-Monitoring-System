package domain

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	Status      string
	CreatedAt   string
}

type GetRepoInfoReq struct {
	Repo  string
	Owner string
}
