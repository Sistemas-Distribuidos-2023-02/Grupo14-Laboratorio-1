docker-central:
	sudo docker compose -f central/docker-compose_central.yml up -d

docker-regional:
	sudo docker compose -f regional/docker-compose_regional.yml up -d

docker-rabbit:
	sudo docker compose -f rabbit/docker-compose_rabbit.yml up -d
