# Audit Server

This service provides an http interface to receive audit logs from external systems such as kong api calls shown [here](https://github.com/giovanibrioni/kong-k8s-monitoring-tools#configure-autid-log) and stores this information in some types of datastore.


## Dependencies

- [Docker](https://docs.docker.com/engine/install/)


## Choose the Storage type
Its possible to choose one of the following storage types: stdout, redis and kafka.
Switch the environment variable STORAGE_TYPE choosing one of the above types on docker-compose.yml file

```bash
  audit-server:
    environment:
      - STORAGE_TYPE=kafka
```

## STORAGE_TYPE=redis
The redis server is configured as following
```bash
  audit-server:
    environment:
      - STORAGE_TYPE=redis
      - REDIS_URL=redis-server:6379
      - REDIS_PASSWORD=
      - REDIS_KEY=audit_logs
```
Starting Redis server
```bash
docker-compose up -d redis-server
```
Than start audit server
```bash
docker-compose up -d audit-server
```

## STORAGE_TYPE=kafka
The kafka broker is configured as following
```bash
  audit-server:
    environment:
      - STORAGE_TYPE=kafka
      - KAFKA_URL=kafka-broker:29092
      - KAFKA_TOPIC=audit_logs
```
Starting Zookeeper
```bash
docker-compose up -d zookeeper
```
Than start Kafka
```bash
docker-compose up -d kafka-broker
```
Than start audit server
```bash
docker-compose up -d audit-server
```

## Generate load
To measure results with more precision, the continer cpu and memory are limited as following:
```bash
mem_limit: 1024m
cpus: 1
```
Running load test
```bash
sh ./k6-loadtest/run-k6-docker.sh
```

## Cleanup

```bash
docker-compose down
```

## EXTRA: Deploy Audit Server on kubernetes
Change the environment variable STORAGE_TYPE and storage server configuration like URL and password
```bash
kubectl apply -f k8s --recursive
```

## Author

Managed by [Giovani Brioni Nunes](https://github.com/giovanibrioni)
