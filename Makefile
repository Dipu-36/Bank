.PHONY: build run test start stop clean db-shell migrate-up 
build:
	@go build -o bin/Bank
run: build
	@./bin/Bank
test: 
	@go test -v ./...
#Docker operations
up:
	@docker-compose up -d --build
	@echo "Containers built & started in background"
down:
	@docker-compose down
	@echo "Contianers stopped and removed"
start:
	@docker start bank-gobank-app-1 bank-gobank-db-1
	@echo "Starting containers...."
stop:
	@docker stop bank-gobank-app-1 bank-gobank-db-1 
	@echo "Continaers have been stopped"
logs:
	@docker-compose logs -f
clean: stop
	@rm -rf bin/
	@docker system prune -f
	@echo "Cleaned all build artifacts and Docker resources"

#Databaase operations
db-shell: 
	@docker exec gobank-app bin/Bank migrate-up

