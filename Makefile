.DEFAULT_GOAL := build

BUILD_DIR := ./cmd/app/build/
PROJECT := ./cmd/app

.PHONY:fmt
fmt:
	go fmt $(PROJECT)/...

.PHONY:vet
vet: fmt
	go vet $(PROJECT)/...

.PHONE:lint
lint: fmt
	go lint $(PROJECT)/...

.PHONY:build
build: vet
	go build -o $(BUILD_DIR) $(PROJECT)/...

.PHONY:run
run: vet
	go run $(PROJECT)/...
