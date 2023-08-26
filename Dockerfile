# 1.Save modules
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /module
# RUN go mod download

# 2.Build project
FROM golang:alpine as builder
WORKDIR /app
# COPY --from=modules /go/pkg /go/pkg
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ./cmd/app/build/ ./cmd/app/main.go

# 3.Execute
FROM scratch
COPY --from=builder /app/config /config
MAINTAINER D3L1K4T3S zhelagin.egor@yandex.ru
LABEL authors="zhelagin.egor"
LABEL version="local"
# CMD ["./cmd/app/build/main"]
