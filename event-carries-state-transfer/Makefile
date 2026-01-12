# Available targets:
# install-tools    - Install required development tools (buf, protoc, etc.)
# generate         - Generate protobuf code and other assets
# setup-hooks      - Setup git hooks for commit message validation
# check-imports    - Check Go import organization (bash script - slower)
# check-imports-go - Check Go import organization (Go script - faster, preferred)
# fix-imports      - Auto-fix import organization issues
# start-infra      - Start NATS and PostgreSQL containers
# start-nats       - Start only NATS container
# start-postgres   - Start only PostgreSQL container
# stop-infra       - Stop NATS and PostgreSQL containers

PHONY := install-tools generate setup-hooks check-imports check-imports-go fix-imports start-infra start-nats start-postgres stop-infra
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

check-imports:
	@echo checking import organization
	@chmod +x scripts/check-imports.sh
	@./scripts/check-imports.sh .
	@echo ✅ import check completed

check-imports-go:
	@echo checking import organization with Go tool
	@go run scripts/check-imports.go .
	@echo ✅ import check completed

fix-imports:
	@echo fixing import organization
	@goimports -w -local eda-in-golang .
	@echo ✅ imports fixed

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
