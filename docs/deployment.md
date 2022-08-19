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
