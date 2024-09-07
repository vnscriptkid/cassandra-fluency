up:
	docker compose up

down:
	docker compose down --volumes --remove-orphans

status:
	docker compose exec -it cassandra-node1 nodetool status

cql:
	docker compose exec -it cassandra-node1 cqlsh

logs2:
	docker compose logs cassandra-node2

logs1:
	docker compose logs cassandra-node1

logs3:
	docker compose logs cassandra-node3
