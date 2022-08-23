WSP_TOOLS_NAME ?= stalactite/wsp-tools
WSP_TOOLS_VER ?= 0.02

.PHONY: build-server build-client run-test-server

build-server:
	go build ./cmd/wsp_server

build-client:
	go build ./cmd/wsp_client

build-local-client:
	go build ./cmd/wsp_local_client

run-test-server:
	go run ./examples/test_api/main.go

docker-build:
	docker build ./docker/ -t ${WSP_TOOLS_NAME}:${WSP_TOOLS_VER}
