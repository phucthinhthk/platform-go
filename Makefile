.PHONY: gen db-reset migrate-up run tools

build:
	docker compose build

build-nc:
	docker compose build --no-cache

up:
	@make set-up-githooks
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

ps:
	docker compose ps

set-up-githooks:
	git config --local core.hooksPath .githooks
	chmod -R +x .githooks/

# Generate backend + frontend (Alias for codegen)
gen:
	@make codegen

# Install requirement tools
tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/air-verse/air@v1.61.1
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

db-reset:
	docker compose exec -T platform-db mysql -uroot -proot -e "DROP DATABASE IF EXISTS platform_db; CREATE DATABASE platform_db;"

migrate-up:
	migrate -path src/db/migrations -database "mysql://root:root@tcp(localhost:3309)/platform_db" up

# run:
# 	air -c .air.toml

air:
	docker compose exec platform-api air -c .air.toml

lint:
	docker compose exec platform-api golangci-lint run

wire:
	docker compose exec platform-api wire gen ./internal/di/wire.go
	docker compose exec platform-api wire gen -output_file_prefix=test_ ./internal/di/wire-test.go

# OpenAPI Bundle
bundle:
	mkdir -p src/generated/api src/generated/openapi/openapi
	docker compose exec platform-node openapi-generator-cli generate -c platform-openapitools.json

codegen:
	@make bundle
	docker compose exec platform-api oapi-codegen -config config/server.yaml generated/openapi/openapi/openapi.yaml
	docker compose exec platform-api oapi-codegen -config config/models.yaml generated/openapi/openapi/openapi.yaml
# 	@make codegen-ts

codegen-ts:
	docker compose exec platform-node npx openapi2aspida -i=./generated/openapi/openapi/openapi.yaml -o=./generated/backend-api
	@if [ -d "../project-frontend" ]; then \
		echo "Moving backend-api files to ../project-frontend..."; \
		cp -r src/generated/backend-api/@types/ ../project-frontend/src/generated/backend-api/@types/ 2>/dev/null || true; \
	fi
	docker compose exec platform-node rm -rf generated/backend-api

#--------------------- migration start
migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migration name=<migration_name>"; \
		exit 1; \
	fi
	docker compose exec platform-api migrate create -ext sql -dir db/migrations $(name)

migrate-up:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db)/platform_db" up

migrate-down:
	docker compose exec platform-api migrate -path db/migrations -database "mysql://root:root@tcp(platform-db)/platform_db" down -all

migrate-refresh:
	@make migrate-down
	@make migrate-up

migrate-fresh:
	docker compose exec platform-db mysql --user=root --password=root -e 'DROP DATABASE IF EXISTS `platform_db`;'
	docker compose exec platform-db mysql --user=root --password=root -e 'CREATE DATABASE `platform_db`;'
	@make migrate-up

seeder:
	docker compose exec platform-api go run db/seed/seed.go