docker run --name some-mongo -p 27017:27017 -v some-mongo-db:/data/db -e MONGO_INITDB_ROOT_USERNAME="mongoadmin" -e MONGO_INITDB_ROOT_PASSWORD="mymongopass123" -d mongo