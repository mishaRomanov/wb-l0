# Wildberries internship L0 assignment. A service that uses Postgres, Docker, Jetstream and Go
## How to run
### Default PG User, password and DB are set in docker-compose file
Clone the repository to your machine: <br>``git clone https://github.com/mishaRomanov/wb-l0`` <br>
Launch docker compose<br>``docker compose up -d`` <br>
Then run the main application<br> ``go run cmd/sub/main.go``<br>
Run application and then run publishing script <br>``cd cmd/pub && go run main.go``

## How to use 
The only endpoint accepts GET requests such as ``localhost:8080/id/b563feb7b2b84b6test``


