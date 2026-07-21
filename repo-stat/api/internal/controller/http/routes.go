package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
	_ "repo-stat/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, agu *usecase.ApiGatewayUsecase) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping, agu))
	mux.Handle("GET /api/repositories/info", NewGetInfoRepositoryHandler(log, agu))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("POST /subscriptions", NewCreateSubscriptionHandler(log, agu))
	mux.Handle("DELETE /subscriptions/{owner}/{repo}", NewDeleteSubscriptionHandler(log, agu))
	mux.Handle("GET /subscriptions", NewGetSubscriptionsHandler(log, agu))
	mux.Handle("GET /subscriptions/info", NewGetInfoRepositoriesHandler(log, agu))
}
