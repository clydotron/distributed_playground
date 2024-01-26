To build docker image:
docker build -f broker-service.dockerfile -t clydotron/broker-service:0.0.1 .

push it to docker hub:
docker push clydotron/broker-service:0.0.1

notes:
https://docs.docker.com/go/access-tokens/

Creating a docker swarm:
project/swarm.yml -- very similar to docker compose

    docker swarm init

    docker swarm join-token worker
    docker swarm join-token manager

docker stack deploy -c swarm.yml myapp

docker service scale <service name>=<num instances>

// shutdown the swarm
docker stack rm myapp

//

docker swarm leave
