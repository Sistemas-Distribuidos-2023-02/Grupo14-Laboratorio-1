docker-central:
	sudo docker compose -f central/docker-compose.central.yml up -d

docker-regional:
	sudo docker compose -f regional/docker-compose.regional.yml up -d

docker-rabbit:
	sudo docker compose -f rabbit/docker-compose.rabbit.yml up -d