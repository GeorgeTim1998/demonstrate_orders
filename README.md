# demonstrate_orders
## Prerequisite
- go 1.21
- psql
- docker 
- nats-streaming docker image sudo docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223
## Set up the project
- run go get .
- run migrations
- start nats-streaming docker 
