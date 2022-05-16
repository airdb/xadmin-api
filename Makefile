VER:='github.com/go-masonry/mortar/mortar.version=v1.2.3'
GIT:='github.com/go-masonry/mortar/mortar.gitCommit=$(shell git rev-parse --short HEAD)'
BUILD_TAG:='github.com/go-masonry/mortar/mortar.buildTag=42'
BUILD_TS:='github.com/go-masonry/mortar/mortar.buildTimestamp=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'

export JAEGER_AGENT_HOST = localhost
export JAEGER_AGENT_PORT = 6831
export JAEGER_SAMPLER_TYPE = const
export JAEGER_SAMPLER_PARAM = 1

run: buf
	@go run -ldflags="-X ${VER} -X ${GIT} -X ${BUILD_TAG} -X ${BUILD_TS}" main.go \
		config config/config.yml

dev-local: buf
	go fmt .
	@CGO_ENABLED=1 go run -ldflags="-X ${VER} -X ${GIT} -X ${BUILD_TAG} -X ${BUILD_TS}" main.go \
		config config/config.yml \
		--additional-files config/config_local.yml

build: buf
	@CGO_ENABLED=1 go build -ldflags="-X ${VER} -X ${GIT} -X ${BUILD_TAG} -X ${BUILD_TS}" main.go

buf:
	@buf generate
	@buf generate --template buf.gen.patch.yaml
	@go install ./cmd/protoc-gen-kit
	@buf generate --template buf.gen.kit.yaml

plugins:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	@go install go.einride.tech/aip/cmd/protoc-gen-go-aip@v0.54.1
	@go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
	@buf generate
	@go install ./cmd/protoc-gen-patch

test:
	@echo "Testing ..."
	@go test -failfast ./...

.PHONY: gen-api test run
