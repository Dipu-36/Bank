.PHONY: build run test start stop clean db-shell migrate-up 
build:
	@go build -o bin/Bank
run: build
	@./bin/Bank
test: 
	@go test -v ./...
#Docker operations
start:
	@docker-compose up -d --build
	@echo "Containers started in background"
stop:
	@docker-compose down
	@echo "Contianers stopped and removed"
logs:
	@docker-compose logs -f
clean: stop
	@rm -rf bin/
	@docker system prune -f
	@echo "Cleaned all build artifacts and Docker resources"

#Databaase operations
db-shell: 
	@docker exec gobank-app bin/Bank migrate-up

