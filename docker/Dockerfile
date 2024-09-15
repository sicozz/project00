FROM golang:1.23 AS builder
ARG TAG="project00-0.0.0"
ENV TAG=${TAG}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/${TAG} ./cmd/server/main.go

# Stage 2: Runner
FROM debian:12-slim AS runner
ARG TAG="project00-0.0.0"
ENV TAG=${TAG}
WORKDIR /
COPY --from=builder /app/hosts.txt /hosts.txt
COPY --from=builder /app/${TAG} /${TAG}
CMD /${TAG}
