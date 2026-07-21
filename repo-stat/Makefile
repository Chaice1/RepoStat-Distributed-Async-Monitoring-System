SHELL := /bin/sh

GOBIN := $(shell go env GOPATH)/bin

PROTOLINT := $(GOBIN)/protolint
GOIMPORTS := $(GOBIN)/goimports
GOLANGCI_LINT := $(GOBIN)/golangci-lint
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc
SQLC:=$(GOBIN)/sqlc


LOCAL_BIN := $(CURDIR)/bin
BUF_BUILD := $(LOCAL_BIN)/buf

PROTOC_GEN_GO_VERSION := v1.32.0
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1
GOLANGCI_LINT_VERSION := v2.11.2



.PHONY: lint fix protobuf protolint gofmtcheck golint tools test

lint: protolint gofmtcheck golint

generate:
	set "PATH=$(LOCAL_BIN);%PATH%" && $(BUF_BUILD) generate


sqlc:
	$(SQLC) generate

protolint:
	$(PROTOLINT) .

gofmtcheck:
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	git diff --exit-code

golint:
	$(GOLANGCI_LINT) run --timeout=2m -E gocritic -v ./...
tools:
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

DB_DSN1 := "postgres://postgres:Ivbln173@localhost:5432/subscriptions?sslmode=disable"
DB_DSN2 := "postgres://postgres:Ivbln173@localhost:5432/repositories?sslmode=disable"

migrate-up:
	migrate -path migrations -database $(DB_DSN1) up
	migrate -path migrations -database $(DB_DSN2) up

migrate-down:
	migrate -path migrations -database $(DB_DSN1) down
	migrate -path migrations -database $(DB_DSN2) down

	 
.tidy:
	go mod tidy

