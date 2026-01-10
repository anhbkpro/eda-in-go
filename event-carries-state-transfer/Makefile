PHONY := install-tools generate setup-hooks start-infra start-nats start-postgres stop-infra
all: install-tools generate setup-hooks

install-tools:
	@echo installing tools
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo ✅ done

generate:
	@echo running code generation
	@go generate ./...
	@echo ✅ done

setup-hooks:
	@echo setting up git hooks
	@echo "Run 'npm install' to install Node.js dependencies and configure git hooks"
	@echo "This sets up commitlint for commit message validation"
	@echo ✅ done

start-infra:
	@echo starting NATS and PostgreSQL
	@docker-compose up -d nats postgres
	@echo ✅ infrastructure started
	@echo "NATS: http://localhost:14222"
	@echo "PostgreSQL: localhost:15432"

start-nats:
	@echo starting NATS
	@docker-compose up -d nats
	@echo ✅ NATS started on port 14222

start-postgres:
	@echo starting PostgreSQL
	@docker-compose up -d postgres
	@echo ✅ PostgreSQL started on port 15432

stop-infra:
	@echo stopping NATS and PostgreSQL
	@docker-compose down nats postgres
	@echo ✅ infrastructure stopped
