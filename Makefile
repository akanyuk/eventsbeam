BUILD_FOLDER := .build
WIN_BINARY_NAME := eventsbeam.exe
LINUX_BINARY_NAME := eventsbeam

.PHONY: all clean windows linux build gen-swagger gen-static

all: build

build: windows linux ## Default: build for windows and linux

windows: gen-static ## Build artifacts for windows
	@echo
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BUILD_FOLDER: $(BUILD_FOLDER)
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BINARY_NAME: $(LINUX_BINARY_NAME)
	@echo
	mkdir -p $(BUILD_FOLDER)/windows
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc  go build -o $(BUILD_FOLDER)/windows/$(WIN_BINARY_NAME) .

linux: gen-static ## Build artifacts for linux
	@echo
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BUILD_FOLDER: $(BUILD_FOLDER)
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BINARY_NAME: $(LINUX_BINARY_NAME)
	@echo
	mkdir -p $(BUILD_FOLDER)/linux
	env GOOS=linux GOARCH=amd64 go build -ldflags -o $(BUILD_FOLDER)/linux/$(LINUX_BINARY_NAME) .

vendor: ## Get dependencies according to go.sum
	env GO111MODULE=auto go mod vendor

test: vendor ## Start unit-tests
	go test ./...

lint: vendor ## Start static code analysis
	hash golangci-lint 2>/dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --timeout=5m

gen-static: gen-swagger ## Generate static resources
	mkdir -p ./web/static/sources
	go run ./generators/static/static.go

gen-swagger: vendor ## Generate swagger.json from CI
	hash swagger 2>/dev/null || go get -u github.com/go-swagger/go-swagger/cmd/swagger
	mkdir -p ./web/static/sources
	go run ./generators/swagger/swagger.go

generate-web-proto: vendor ## Generate pr.go for api generation plugin
	protoc --proto_path=./cmd/protoc-gen-web/proto/web --go_out=./cmd/protoc-gen-web/proto/web web.proto
	cp ./cmd/protoc-gen-web/proto/web/*.proto ./api/

build-protoc-gen-web-windows: vendor ## Build api generation plugin
	mkdir -p ./api
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o protoc-gen-web.exe ./cmd/protoc-gen-web/

generate-api: generate-web-proto build-protoc-gen-web-windows ## Generate API
	protoc -I=cmd/protoc-gen-web/proto/web --proto_path=./api --web_out=./api --go_out=./api api.proto
	rm ./protoc-gen-web.exe || true

clean: ## Remove vendor and artifacts
	rm -rf vendor
	rm -rf $(BUILD_FOLDER)/linux
	rm -rf $(BUILD_FOLDER)/windows

	rm ./api/web.proto || true
	rm ./api/descriptor.proto || true

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 
