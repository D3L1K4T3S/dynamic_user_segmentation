.DEFAULT_GOAL := build

BUILD_DIR := ./cmd/app/build/
PROJECT := .

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
	go build $(PROJECT)/... -o $(BUILD_DIR)

.PHONY:run
run: vet
	go run $(PROJECT)/...
