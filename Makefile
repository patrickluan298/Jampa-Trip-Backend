docker-dev-up: # Inicia os containers em background
	docker-compose -f deployments/docker-compose.yaml up -d

docker-dev-build: # Faz o build e inicia os containers
	docker-compose -f deployments/docker-compose.yaml up --build -d

docker-dev-logs: # Exibe logs em tempo real
	docker-compose -f deployments/docker-compose.yaml logs -f --tail=100

# Faz o build + inicia containers + exibe logs
docker-dev-build-logs: docker-dev-build docker-dev-logs

docker-dev-stop: # Para containers sem removê-los
	docker-compose -f deployments/docker-compose.yaml stop

docker-dev-down: # Para e remove containers/volumes
	docker-compose -f deployments/docker-compose.yaml down