
check-config:
	docker compose --env-file config.env config

docker-compose-run:
	docker compose --env-file config.env -p walletapp up -d

docker-compose-stop:
	docker compose --env-file config.env -p walletapp down

docker-compose-build:
	docker compose --env-file config.env -p walletapp build
