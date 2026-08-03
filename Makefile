include .env
export

ifeq ($(OS),Windows_NT)
    PROJECT_ROOT := $(CURDIR)
else
    PROJECT_ROOT := $(shell pwd)
endif
export PROJECT_ROOT

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Clean all environment volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres port-forwarder && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "Files cleanup"; \
	else \
	  echo "Cleanup is closed"; \
	fi

logs-cleanup:
	@read -p "Clean all log files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  rm -rf ${PROJECT_ROOT}/out/logs && \
	  echo "Files cleanup"; \
	else \
	  echo "Cleanup is closed"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Param seq is not exist: " \
  		exit 1; \
  	fi
	@docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
    		echo "Param action is not exist: "; \
    		exit 1; \
    	fi
	@docker compose run --rm --use-aliases todoapp-postgres-migrate \
		-path=//migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy:
	@docker compose up -d --build todoapp

ps:
	@docker compose ps
