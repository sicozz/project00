FROM golang:1.22 AS builder
ARG TAG="project00-0.0.0"
ENV TAG=$TAG
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /$TAG .

# FROM alpine:latest
FROM debian:12-slim AS runner
ARG TAG="project00-0.0.0"
ENV TAG=$TAG
WORKDIR /
COPY --from=builder /$TAG .
# EXPOSE 50051
CMD /$TAG
