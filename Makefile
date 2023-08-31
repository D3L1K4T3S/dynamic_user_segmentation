include .env
export

compose-up: ### Запуск docker compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Остановка docker compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag: ### Сгенинерировать swagger docs
	swag init -g /cmd/app/main.go --parseDependency

test: ### Запуск тестов
	go test -v ./...

cover-html: ### Запуск тестов  покрытия
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: coverage-html

mockgen: ### Генерирования 'моков' для тестирования
	mockgen -source=internal/service/service.go 		-destination=internal/mocks/service/service.go 			   -package=servicemocks
	mockgen -source=internal/repository/repository.go   -destination=internal/mocks/repository/repository.go       -package=repomocks
.PHONY: mockgen

migrate-create:  ### Создание миграции базы данных
	migrate create -ext sql -dir migrations 'dynamic_segments'
.PHONY: migrate-create

migrate-up: ### Сделать миграцию
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### Убрать миграцию
	echo "y" | migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

linter-golangci: ### Проверка линтером golangci
	golangci-lint run
.PHONY: linter-golangci

help: ## Все возможные команды для работы с проектом
	@awk 'BEGIN {FS = ":.*##"; printf "\nКоманды:\n  make \033[36m<команда>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help