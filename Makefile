NAME = zk-axon
IMAGE_NAME = zk-axon
IMAGE_VERSION = 1.0
IMAGE_VERSION_MULTI_ARCH = multiarch

LOCATION ?= us-west1
PROJECT_ID ?= zerok-dev
REPOSITORY ?= zk-axon

BUILDER_NAME = multi-platform-builder
IMAGE_PREFIX := $(LOCATION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/

export GO111MODULE=on
export GOPRIVATE=github.com/zerok-ai/zk-utils-go,github.com/zerok-ai/zk-rawdata-reader

sync:
	go get -v ./...

build: sync
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(NAME) cmd/main.go

docker-build: sync
	CGO_ENABLED=0 GOOS=linux $(ARCH) go build -v -o $(NAME) cmd/main.go
	docker build --no-cache -t $(IMAGE_PREFIX)$(IMAGE_NAME):$(IMAGE_VERSION) .

build-multiarch: sync
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-amd64 cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/$(NAME)-arm64 cmd/main.go
	#Adding remove here again to account for the case when buildx was not removed in previous run.
	docker buildx rm ${BUILDER_NAME} || true
	docker buildx create --use --platform=linux/arm64,linux/amd64 --name ${BUILDER_NAME}
	docker buildx build --platform=linux/arm64,linux/amd64 --push \
	--tag $(IMAGE_PREFIX)$(IMAGE_NAME):$(IMAGE_VERSION) .
	docker buildx rm ${BUILDER_NAME}

docker-build-gke: IMAGE_PREFIX := $(LOCATION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/
docker-build-gke: ARCH := GOARCH=amd64
docker-build-gke: docker-build

docker-push-gke: IMAGE_PREFIX := $(LOCATION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/
docker-push-gke:
	docker push $(IMAGE_PREFIX)$(IMAGE_NAME):$(IMAGE_VERSION)

docker-build-push-gke: docker-build-gke docker-push-gke

run: build
	go run cmd/main.go -c ./config/config.yaml 2>&1 | grep -v '^(0x'

dev: sync
	nodemon --watch './**/*.go' --signal SIGTERM --exec make run

fmt:
	gofmt -s -w .

test:
	go test ./... -cover

coverage_cli:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage_html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ------- CI-CD ------------
ci-cd-build: sync
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o bin/$(NAME)-amd64 cmd/main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -v -o bin/$(NAME)-arm64 cmd/main.go

