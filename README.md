# demonstrate_orders
## Description
This is a service that receives messages via nats-streaming, caches the orders info and stores the orders in the postgres database. Info about the order can be requested via filling the form at `/` with uid of the order.
## Prerequisite
- go `1.21`
- psql `16`
- docker `26.1.3`
- nats-streaming `docker image sudo docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223`
## Set up the project
- run go get .
- run migrations `./bin/goose -dir db/migrations postgres "postgres://postgres:PWD@localhost:5432/demonstrate_orders?sslmode=disable"  up`
- start nats-streaming docker `docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223`
- Run the app with `go run main.go` or build it and run `go build -o main . && ./main`
## Tests
- Tests are implemented with the help from docker.
- To run tests run `docker compose --build`.
