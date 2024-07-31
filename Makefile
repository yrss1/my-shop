SERVICE_NAMES = product-service-app payment-service-app order-service-app user-service-app api-gateway-service-app

build:
	@echo "Building Docker images..."
	@for service in $(SERVICE_NAMES); do \
		docker build -t $$service ./microservices/$$service; \
	done
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

.PHONY: build up down restart


#docker-compose up --build
#test:
#	@echo "Running tests..."

