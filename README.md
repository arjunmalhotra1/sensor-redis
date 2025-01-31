# sensor-redis
## steps to run
1. `docker-compose up -d`
2. `go run main.go`
## Example curl post request
```curl --location 'localhost:8086/message' \
--header 'Content-Type: application/json' \
--data '{
    "device_id":"1234",
    "device_type":"weather",
    "temp":3.00

}'