.PHONY: dev dev-backend dev-frontend build-frontend build-all docker-build docker-run test lint

dev:
	@echo "Starting development servers..."
	@make -j 2 dev-backend dev-frontend

dev-backend:
	GIN_MODE=debug go run main.go

dev-frontend:
	cd ui && npm run dev

create-migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir database/migrations -seq $$name

build-frontend:
	cd ui && npm run build

build-all: build-frontend
	GIN_MODE=release go build -o main .

docker-build:
	docker build -t ghost-listmonk-connector .

docker-run:
	docker run -p 8808:8808 ghost-listmonk-connector

repopack:
	npx repopack --ignore "*.txt,*.db,go.sum,ui_backup,**/*.txt,**/ui/*.tsx,**/ui/*.ts,**/.next,**/node_modules,**/pnpm-lock.yaml" 

test:
	go test ./...
	cd ui && npm test

lint:
	go vet ./...
	cd ui && npm run lint