BUILD_FOLDER := .build
WIN_BINARY_NAME := eventsbeam.exe
LINUX_BINARY_NAME := eventsbeam

.PHONY: all clean windows linux build etc-windows etc-linux gen-swagger

all: build

build: windows linux ## Default: build for windows and linux

windows: vendor etc-windows ## Build artifacts for windows
	@echo
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BUILD_FOLDER: $(BUILD_FOLDER)
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BINARY_NAME: $(LINUX_BINARY_NAME)
	@echo
	mkdir -p $(BUILD_FOLDER)/windows
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc  go build -o $(BUILD_FOLDER)/windows/$(WIN_BINARY_NAME) .

etc-windows:
	mkdir -p $(BUILD_FOLDER)/windows
	cp -r etc/* $(BUILD_FOLDER)/windows/

linux: vendor etc-linux ## Build artifacts for linux
	@echo
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BUILD_FOLDER: $(BUILD_FOLDER)
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" BINARY_NAME: $(LINUX_BINARY_NAME)
	@echo
	mkdir -p $(BUILD_FOLDER)/linux
	env GOOS=linux GOARCH=amd64 go build -ldflags -o $(BUILD_FOLDER)/linux/$(LINUX_BINARY_NAME) .

etc-linux:
	mkdir -p $(BUILD_FOLDER)/linux
	cp -r etc/* $(BUILD_FOLDER)/linux/

vendor: ## Get dependencies according to go.sum
	env GO111MODULE=auto go mod vendor

test: vendor ## Start unit-tests
	go test ./...

lint: vendor ## Start static code analysis
	hash golangci-lint 2>/dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --timeout=5m

gen-swagger: vendor ## Generate swagger.json from CI
	mkdir -p ./web/static/sources
	go run ./generators/web_resources/web_resources.go -stage swagger -metadata ./metadata.yml -swagger-work-dir ./web -swagger-output-file ./web/static/sources/swagger.json

docker-run: linux ## Run service in docker
	@echo
	@printf "\033[35m%-18s\033[33m %s\033[0m\n" SED_DOCKER_NETWORK: $(SED_DOCKER_NETWORK)
	@echo
	@mkdir -p $(BUILD_FOLDER)/docker-run
	@cp -r $(BUILD_FOLDER)/linux/* $(BUILD_FOLDER)/docker-run
	@cp -r docker/docker-run/* $(BUILD_FOLDER)/docker-run

	docker container rm --force $(BINARY_NAME) || true
	docker build --quiet --force-rm --tag $(BINARY_NAME):$(BINARY_VERSION) $(BUILD_FOLDER)/docker-run
	docker run --detach --rm --name $(BINARY_NAME) --network $(SED_DOCKER_NETWORK) $(BINARY_NAME):$(BINARY_VERSION)

clean: ## Remove vendor and artifacts
	rm -rf vendor
	rm -rf $(BUILD_FOLDER)/linux
	rm -rf $(BUILD_FOLDER)/windows

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 
