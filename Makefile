API_PROTO_FILES=$(shell find api -name *.proto)
ENV="local"
BIN_PATH=bin/
CMD_PATH=cmd/
DEPLOY_PATH=deploy/
SERVICE_LIST = gateway logic msg transfer

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	#go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	#go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

.PHONY: wire
wire:
	cd cmd/logic && wire
	cd cmd/gateway && wire
	cd cmd/msg && wire
	cd cmd/transfer && wire

.PHONY: proto
proto:
	protoc --proto_path=./api/proto \
           --proto_path=./third_party \
           --go_out=paths=source_relative:./api/pb \
           --go-grpc_out=paths=source_relative:./api/pb \
           --go-errors_out=paths=source_relative:./api/pb \
           --validate_out=paths=source_relative,lang=go:./api/pb \
           $(API_PROTO_FILES)

service-build:
	@$(foreach item,$(SERVICE_LIST),CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BIN_PATH)$(item) ./$(CMD_PATH)$(item) |) echo "go rpc build running"
docker-build:
	@$(foreach item,$(SERVICE_LIST),docker build -f $(DEPLOY_PATH)Dockerfile --build-arg ENV=$(ENV) --build-arg FILENAME=$(item) -t $(item) . |) echo "docker api build running"
.PHONY: dc
dc:
	docker-compose -f deploy/docker-compose.yaml up -d

.PHONY: stop
stop:
	docker-compose -f deploy/docker-compose.yaml stop

.PHONY: start
start: service-build docker-build dc

.PHONY: dall
dall:
	docker image prune -f
	docker system prune -a


