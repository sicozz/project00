APP				?= project00
APP_VERSION		?= 0.0.1
TAG				?= $(APP)-$(APP_VERSION)
SERVER_ENTRY	?= cmd/$(APP)/main.go
SERVER_BIN		?= bin/$(TAG)
DOCKER_TAG		?= $(APP):$(APP_VERSION)
DOCKER_CONF		?= ./docker
API_VERSION		?= 0.0
PROTO_DIR		?= api/v$(API_VERSION)

# DEBUG

debug:
	dlv debug ./cmd/$(APP)

# BUILD
build: build-proto build-server

build-proto:
	@ protoc --proto_path=$(PROTO_DIR) \
       --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
       --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
       $(PROTO_DIR)/proto00.proto

build-server:
	@ go build -o $(SERVER_BIN) $(SERVER_ENTRY)

install-dev-deps:
	@ go install github.com/go-delve/delve/cmd/dlv@latest
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# DOCKER
docker-build:
	@ docker build \
		-f $(DOCKER_CONF)/Dockerfile \
		--build-arg TAG=$(TAG) \
		-t $(DOCKER_TAG) .

docker-clean:
	@ docker rmi $(DOCKER_TAG)

# CLEAN
clean:
	@ rm -r $(SERVER_BIN)

# DISTENV
distenv-up:
	@ docker compose \
		--env-file $(DOCKER_CONF)/.env \
		-f $(DOCKER_CONF)/compose.yaml up \
		-d

distenv-down:
	@ docker compose \
		-f $(DOCKER_CONF)/compose.yaml down

distenv-monitor:
	docker compose \
		-f $(DOCKER_CONF)/compose.yaml logs \
		-f --no-color \
		| fzf --delimiter=' ' --nth=1

distenv-connect:
	NODE=$(NODE) docker exec -it $(NODE) /bin/bash

distenv-run: distenv-down build-server docker-build distenv-up

# RUN
run:
	./$(SERVER_BIN)

# LINT
lint:
	@ golines --max-len=80 -w .
