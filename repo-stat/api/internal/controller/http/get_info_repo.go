package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	api_controller_errors "repo-stat/api/internal/controller/errors"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"
)

// GetInfoRepositoryHandler
// @Summary      Получить инфо о репозитории
// @Description  Запрашивает данные о звездах, форках и коммитах
// @Tags         repository
// @Produce      json
// @Param 		 url	 query     string  true	"Полная ссылка на репозиторий (например, https://github.com/golang/go)"
// @Success      200     {object}  dto.RepoInfo
// @Failure      400     {object}  map[string]string
// @Failure		 404     {object}  map[string]string
// @Failure 	 500	 {object}  map[string]string
// @Router       /api/repositories/info [get]
func NewGetInfoRepositoryHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url_str := r.URL.Query().Get("url")

		if url_str == "" {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, domain.ErrBadRequest, log, "get url string")
			return
		}

		parsed_url, err := url.Parse(url_str)

		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, domain.ErrBadRequest, log, "parsing url string")
			return
		}

		path := strings.Trim(parsed_url.Path, "/")
		path_slice := strings.Split(path, "/")

		if len(path_slice) != 2 {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, domain.ErrBadRequest, log, "check size of path_slice")
			return
		}

		repo_info, err := agu.GetInfoRepo(r.Context(), []*domain.GetRepoInfoReq{
			&domain.GetRepoInfoReq{Repo: path_slice[1], Owner: path_slice[0]},
		})

		log.Error("api", "error", err, "repo_info", repo_info)

		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "GetInfoRepo")
			return
		}

		resp := dto.RepoInfo{
			FullName:    repo_info[0].FullName,
			Description: repo_info[0].Description,
			Forks:       repo_info[0].Forks,
			Stargazers:  repo_info[0].Stargazers,
			CreatedAt:   repo_info[0].CreatedAt,
			Status:      repo_info[0].Status,
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to write RepoInfo", "error", err)
		}
	}
}

// GetInfoRepositoriesHandler
// @Summary      Получить инфо о репозиториях, на которые подписан пользователь
// @Description  Запрашивает данные о звездах, форках и коммитах
// @Tags         repository
// @Produce      json
// @Success      200     {object}  dto.GetRepositoriesInfoResponse
// @Failure      400     {object}  map[string]string
// @Failure		 404     {object}  map[string]string
// @Failure 	 500	 {object}  map[string]string
// @Router       /subscriptions/info [get]
func NewGetInfoRepositoriesHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		subscriptions, err := agu.GetSubscriptions(r.Context())

		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "GetSubscriptions")
			return
		}

		req := make([]*domain.GetRepoInfoReq, len(subscriptions))

		for i := range subscriptions {
			req[i] = &domain.GetRepoInfoReq{
				Repo:  subscriptions[i].RepoName,
				Owner: subscriptions[i].OwnerName,
			}
		}

		repositories, err := agu.GetInfoRepo(r.Context(), req)

		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "GetInfoRepo")
			return
		}

		resp := dto.GetRepositoriesInfoResponse{Repositories: make([]*dto.RepoInfo, len(subscriptions))}

		for i, item := range repositories {

			resp.Repositories[i] = &dto.RepoInfo{
				FullName:    item.FullName,
				Description: item.Description,
				Forks:       item.Forks,
				Stargazers:  item.Stargazers,
				CreatedAt:   item.CreatedAt,
				Status:      item.Status,
			}

		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to write RepositoriesInfo", "error", err)
		}

	}
}
