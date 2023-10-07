docker-central:
	sudo docker compose -f Central/docker-compose_central.yml up -d

docker-regional:
	sudo docker compose -f Regional/docker-compose_regional.yml up -d

docker-rabbit:
	sudo docker compose -f rabbit/docker-compose_rabbit.yml up -d
