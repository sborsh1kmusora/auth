include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(MIGRATION_DSN)

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.25.0

migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/pressly/goose/v3/cmd/goose

generate:
	make generate-user-api
	make generate-auth-api
	make generate-access-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	api/user_v1/user.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	api/auth_v1/auth.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	api/access_v1/access.proto

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/sborsh1kmusora/auth/internal/service/...,github.com/sborsh1kmusora/auth/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/sborsh1kmusora/auth/internal/service/...,github.com/sborsh1kmusora/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore
