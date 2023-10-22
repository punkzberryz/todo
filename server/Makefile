include .env

db-login:
	docker exec -it ${DB_DOCKER_CONTAINER} bash && psql && \connect ${DB_NAME}
	psql
	\connect ${DB_NAME}
postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:12-alpine

createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

dropdb:
	docker exec -it ${DB_DOCKER_CONTAINER} dropdb ${DB_NAME}

dropAndCreate:
	docker exec -it ${DB_DOCKER_CONTAINER} dropdb ${DB_NAME}
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}
	migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" -verbose up
	sqlc generate

migrate_create:
	migrate create -ext sql -dir db/migration -seq init_schema

migrate_up:
	migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" -verbose down

sqlc:
	sqlc generate
	
sqlc_windows:
	docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate

test:
	cp .env app.env
	go test -v -cover ./...

server:		
	cp .env app.env
	go run main.go

stop_containers:
	@echo "Stopping other docker container"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "no containers running..."; \
	fi
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Papai-hospital/papai-server-go/db/sqlc Store
git-pull:
	git fetch && git pull