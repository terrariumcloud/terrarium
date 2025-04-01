# Deployment


## Local deployment

For local deployment use provided docker-compose file. To build image just run:

```bash
docker-compose build
```

This will create an image named `terrarium` with dev `tag`.

Before running the containers, first check `.env` file and provide missing values.

```bash
docker-compose up -d
```

Terrarium UI: http://localhost:50006/

Jaeger: http://localhost:16686/


## Local deployment with localstack

1. Run localstack

        docker run \
          --rm -it \
          -p 127.0.0.1:4566:4566 \
          -p 127.0.0.1:4510-4559:4510-4559 \
          -v /var/run/docker.sock:/var/run/docker.sock \
          --network terrarium \
          --name localstack \
          localstack/localstack

2. Build and run:

        docker compose build

        export USE_LOCALSTACK=true
        docker compose up

