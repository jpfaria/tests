# Developer

## Run Nats

## single

```bash
docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats -DV
```

## cluster

```bash
docker network create nats

docker run -d --name nats-seed --network nats \
    --rm \
    -p 4222:4222 -p 8222:8222 \
    nats \
    --http_port 8222 --cluster_name NATS --cluster nats://0.0.0.0:6222

docker run -d --name nats-1 --network nats \
    --rm \
    nats \
    --cluster_name NATS --cluster nats://0.0.0.0:6222 --routes nats://ruser:T0pS3cr3t@nats-seed:6222

docker run -d --name nats-2 --network nats \
    --rm \
    nats \
    --cluster_name NATS --cluster nats://0.0.0.0:6222 --routes nats://ruser:T0pS3cr3t@nats-seed:6222

docker run -d --name nats-3 --network nats \
    --rm \
    nats \
    --cluster_name NATS --cluster nats://0.0.0.0:6222 --routes nats://ruser:T0pS3cr3t@nats-seed:6222
```