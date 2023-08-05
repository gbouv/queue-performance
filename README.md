Queue performance scripts
=========================

Run postgres container:
```
docker volume create pg-data
docker run --rm -it --name database -p 5432:5432 -e POSTGRES_DB=postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e PGDATA=/data -v pg-data:/data --memory 512m --cpus 0.1 postgres
```

Run postgres container:
```
docker volume create pg-data
docker run --rm -it --name database -p 5432:5432 -e POSTGRES_DB=postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e PGDATA=/data -v pg-data:/data --memory 512m --cpus 0.1 postgres
```

Run redis container:
```
docker volume create redis-data
docker run --rm -it --name redis -p 6379:6379 -v redis-data:/data --memory 512m --cpus 0.1 -d redis
```

Build and run loader container:
```
docker build -t loader:latest -f ./loader/Dockerfile ./loader

for i in {1..4}; do docker run --rm --name "loader-$i" -e DB_HOST=$DB_HOST -d loader:latest; done
```

Build and run consumer container:
```
docker build -t consumer:latest -f ./consumer/Dockerfile ./consumer

for i in {1..8}; do docker run --rm --name "consumer-$i" -e DB_HOST=$DB_HOST -d consumer:latest; done
```
