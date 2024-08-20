APP				?= project00
APP_VERSION		?= 0.0.1
TAG				?= $(APP)-$(APP_VERSION)
CLIENT_TAG		?= $(APP)_client-$(APP_VERSION)
SERVER_ENTRY	?= cmd/server/main.go
CLIENT_ENTRY	?= cmd/client/main.go
SERVER_BIN		?= bin/$(TAG)
CLIENT_BIN		?= bin/$(CLIENT_TAG)
DOCKER_TAG		?= $(APP):$(APP_VERSION)
DOCKER_CONF		?= ./docker
API_VERSION		?= 0.0
PROTO_DIR		?= api/v$(API_VERSION)

# BUILD
build: build-proto build-server build-client

build-proto:
	@ protoc --proto_path=$(PROTO_DIR) \
       --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
       --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
       $(PROTO_DIR)/proto00.proto

build-server:
	@ go build -o $(SERVER_BIN) $(SERVER_ENTRY)

build-client:
	@ go build -o $(CLIENT_BIN) $(CLIENT_ENTRY)

install-dev-deps:
	@ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	@ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

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
	@ rm -r $(SERVER_BIN) $(CLIENT_BIN)

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
