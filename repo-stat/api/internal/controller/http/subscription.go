package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	api_controller_errors "repo-stat/api/internal/controller/errors"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// GetSubscriptionsHandler
// @Summary      Получить инфо о подписках пользователя
// @Description  берёт из бд инфо о подписках
// @Tags         subscription
// @Produce      json
// @Success      200     {object}  dto.GetSubscriptionsResponse
// @Failure      400     {object}  map[string]string
// @Failure		 404     {object}  map[string]string
// @Failure 	 500	 {object}  map[string]string
// @Router       /subscriptions [get]
func NewGetSubscriptionsHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := agu.GetSubscriptions(r.Context())

		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "GetSubscriptions")
			return
		}
		response := dto.GetSubscriptionsResponse{
			Subscriptions: make([]*dto.Subscription, len(subscriptions)),
		}

		for i, item := range subscriptions {
			response.Subscriptions[i] = &dto.Subscription{
				Id:        item.Id,
				RepoName:  item.RepoName,
				OwnerName: item.OwnerName,
				CreatedAt: item.CreatedAt,
			}
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write deletesubscription response", "action", "encode response", "error", err)
		}

	}
}

// DeleteSubscriptionHandler
// @Summary      Отписаться от репозитория
// @Description  Удаляет репозиторий из списка отслеживаемых
// @Tags         subscription
// @Param        owner path      string  true  "Владелец репозитория"
// @Param        repo  path      string  true  "Название репозитория"
// @Success      200   {object}  map[string]string "message: deleted successfully"
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /subscriptions/{owner}/{repo} [delete]
func NewDeleteSubscriptionHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")

		err := agu.DeleteSubscription(r.Context(), owner, repo)
		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "DeleteSubcription")
			return
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{"message": "deleted successfully"}); err != nil {
			log.Error("failed to write deletesubscription response", "action", "encode response", "error", err)
		}
	}
}

// CreateSubscriptionHandler
// @Summary      Подписаться на репозиторий
// @Description  Добавляет репозиторий в список отслеживаемых
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateSubscriptionRequest  true  "Данные подписки"
// @Success      200   {object}  map[string]string "message: create subscription successfully"
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /subscriptions [post]
func NewCreateSubscriptionHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody := dto.CreateSubscriptionRequest{}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, domain.ErrBadRequest, log, "Parse CreateSubscription req body")
			return
		}
		err := agu.CreateSubscription(r.Context(), reqBody.OwnerName, reqBody.RepoName)
		if err != nil {
			api_controller_errors.HandleErrorsFromDomainToHTTP(w, err, log, "CreateSubscription")
			return
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{"message": "create subscription successfully"}); err != nil {
			log.Error("failed to write createsubscription response", "action", "encode response", "error", err)
		}
	}
}
