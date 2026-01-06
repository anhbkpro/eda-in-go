PHONY := install-tools generate setup-hooks
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
