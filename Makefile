APP_NAME=gophkeeper

BUILD_DIR := build
BIN_DIR := bin
DB_HOST=192.168.0.105
DB_USER=gophkeeper
DB_NAME=gophkeeper
DB_PASS=gophkeeper
DB_PORT=5432
DB_STRING="postgres://$(DB_NAME):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
DB_MIGRATIONS_PATH="./migrations"

.DEFAULT_GOAL := help

### protoc
PROTOC_VERSION := $(shell curl -s https://api.github.com/repos/protocolbuffers/protobuf/releases/latest | jq -r '.tag_name')
ifeq ($(PROTOC_VERSION),null) # в случае если запрос нарвется на rate limit
	PROTOC_VERSION := 33.0
else
	PROTOC_VERSION := $(shell echo $(PROTOC_VERSION) | cut -c 2-)
endif
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
PROTOC_DOWNLOAD_LINK := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
PROTOC_BIN := $(BUILD_DIR)/protoc-$(PROTOC_VERSION)
$(PROTOC_BIN):
	@mkdir -p $(BUILD_DIR)
	@wget -q --show-progress -O $(BUILD_DIR)/$(PROTOC_ZIP) $(PROTOC_DOWNLOAD_LINK)
	@unzip -o -q -j $(BUILD_DIR)/$(PROTOC_ZIP) bin/protoc -d $(BUILD_DIR)
	@mv $(BUILD_DIR)/protoc $(BUILD_DIR)/protoc-$(PROTOC_VERSION)
	@rm $(BUILD_DIR)/$(PROTOC_ZIP)
###

define install-tool
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && \
		go mod init temp || true && \
		go get $(1) && \
		go build $(firstword $(subst @, ,$(1))) && \
		rm go.mod go.sum
endef

### protoc-gen-go
PROTOC_GEN_GO := google.golang.org/protobuf/cmd/protoc-gen-go
PROTOC_GEN_GO_BIN := $(BUILD_DIR)/protoc-gen-go
$(PROTOC_GEN_GO_BIN):
	$(call install-tool, $(PROTOC_GEN_GO))
###

### protoc-gen-go-grpc
PROTOC_GEN_GO_GRPC := google.golang.org/grpc/cmd/protoc-gen-go-grpc
PROTOC_GEN_GO_GRPC_BIN := $(BUILD_DIR)/protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC_BIN):
	$(call install-tool, $(PROTOC_GEN_GO_GRPC))
###

### protoc-gen-openapiv2
PROTOC_GEN_OPENAPIV2 := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
PROTOC_GEN_OPENAPIV2_BIN := $(BUILD_DIR)/protoc-gen-openapiv2
$(PROTOC_GEN_OPENAPIV2_BIN):
	$(call install-tool, $(PROTOC_GEN_OPENAPIV2))
###

### protoc-gen-grpc-gateway
PROTOC_GEN_GRPC_GATEWAY := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
PROTOC_GEN_GRPC_GATEWAY_BIN := $(BUILD_DIR)/protoc-gen-grpc-gateway
$(PROTOC_GEN_GRPC_GATEWAY_BIN):
	$(call install-tool, $(PROTOC_GEN_GRPC_GATEWAY)@v2.27.4)
###

### protoc-gen-doc
PROTOC_GEN_DOC := github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
PROTOC_GEN_DOC_BIN := $(BUILD_DIR)/protoc-gen-doc
$(PROTOC_GEN_DOC_BIN):
	$(call install-tool, $(PROTOC_GEN_DOC)@v1.5.1)
###

.PHONY: install-proto-tools
install-proto-tools: $(PROTOC_BIN) $(PROTOC_GEN_GO_BIN) $(PROTOC_GEN_GO_GRPC_BIN) $(PROTOC_GEN_OPENAPIV2_BIN) $(PROTOC_GEN_GRPC_GATEWAY_BIN) $(PROTOC_GEN_DOC_BIN)

PROTO_DIR=proto
PROTO_V1_DIR=$(PROTO_DIR)/$(APP_NAME)/v1
SWAGGER_OUT_DIR=api
.PHONY: proto
proto: install-proto-tools
	@mkdir docs -p
	@$(PROTOC_BIN) -I$(PROTO_V1_DIR) -I$(PROTO_DIR) \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO_BIN) \
		--go_out=$(PROTO_V1_DIR) \
		--go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC_BIN) \
		--go-grpc_out=$(PROTO_V1_DIR) \
		--go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GRPC_GATEWAY_BIN) \
		--grpc-gateway_out=$(PROTO_V1_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(PROTOC_GEN_OPENAPIV2_BIN) \
		--openapiv2_out=json_names_for_fields=false:$(SWAGGER_OUT_DIR) \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt allow_merge=true \
		$(PROTO_V1_DIR)/$(APP_NAME).proto

BINARY_NAME ?= gophkeeper-cli
BUILD_VERSION ?= v1.0.0
BUILD_DATE ?= 2026-03-04
GRPC_SERVER_ADDR ?= :8080
LDFLAGS = -X main.buildVersion=$(BUILD_VERSION) -X main.buildDate=$(BUILD_DATE)
.PHONY: build-cli
build-cli:
	@echo "Building Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./cmd/client/$(BINARY_NAME)-linux   ./internal/client

	@echo "Building macOS..."
	GOOS=darwin  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./cmd/client/$(BINARY_NAME)-darwin  ./internal/client

	@echo "Building Windows..."
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./cmd/client/$(BINARY_NAME)-win.exe ./internal/client

.PHONY: build
build:
	go build -o ./cmd/server/gophkeeper ./cmd/server

.PHONY: migrate-up
migrate-up:
	migrate \
	-path $(DB_MIGRATIONS_PATH) \
	-database $(DB_STRING) up

.PHONY: migrate-down
migrate-down:
	migrate \
	-path $(DB_MIGRATIONS_PATH) \
	-database $(DB_STRING) down


.PHONY: migrate-create
migrate-create:
ifdef NAME
	migrate create \
    	-ext sql \
    	-dir $(DB_MIGRATIONS_PATH) \
    	-seq $(NAME)
else
	@echo "Require variable NAME not found"
endif


.PHONY: install-pg-tools
install-pg-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest # golang-migrate CLI

.PHONY: install-mock-tools
install-mock-tools:
	go install github.com/golang/mock/mockgen@latest  # mocks for tests
	go install github.com/golang/mock/gomock@latest

.PHONY: install-all-tools
install-all-tools: install-proto-tools install-mock-tools install-pg-tools

.PHONY: test
test:
	go test -v ./... | { grep -v 'no test files'; true; }

.PHONY: test_cover
test_cover:
	go test -coverprofile=coverage.out ./...
	cat coverage.out | grep -v '/mocks\|/test\|/vendor\|/internal/model\|/proto' > coverage.filtered.out
	go tool cover -func=coverage.filtered.out
	rm coverage.out coverage.filtered.out

.PHONY: mock
mock:
	@echo "Generating mock for Storage..."
	mockgen -destination=internal/repository/pgstorage/mocks/pgstorage_mock.go -package=pgstorage -source=internal/service/service.go Storage
	@echo "Generating mock for Service..."
	mockgen -destination=internal/service/mocks/service_mock.go -package=service -source=internal/controller/grpc/grpc.go Service

.PHONY: gofmt
gofmt:
	@gofmt -w ./..

CYAN := \033[36m
BOLD := \033[1m
NO_COLOR := \033[0m
LOGO := "🦫"

.PHONY: help
help:
	@echo "$(LOGO)  $(CYAN)$(BOLD)$(APP_NAME)$(NO_COLOR)"
	@echo ""
	@echo "command           | description"
	@echo "===================================================="
	@echo "install-pg-tools       | install golang-migrate CLI"
	@echo "install-mock-tools     | install mockgen/gomock"
	@echo "install-proto-tools    | download protoc + plugins"
	@echo "install-all-tools      | install proto/mock/pg tools"
	@echo "proto                  | generate proto files"
	@echo "migrate-up             | apply DB migrations"
	@echo "migrate-down           | rollback migration"
	@echo "migrate-create         | create migration (NAME=...) "
	@echo "mock                   | generate Storage/Service mocks"
	@echo "test                   | run tests"
	@echo "test_cover             | tests + coverage report"
	@echo "gofmt                  | format all Go files"
	@echo "build-cli              | cross-platform CLI builds"
	@echo "build                  | build server"
