docker run --name natstreaming -d -p 4222:4222 nats-streaming -SD -V -cid whale
docker run --name redis -d -p 6379:6379 redis

