# Wildberries internship L0 assignment. A service that uses Postgres, Docker, Jetstream and Go
## How to run
### Default PG User, password and DB are set in docker-compose file
Clone the repository to your machine: <br>``git clone https://github.com/mishaRomanov/wb-l0`` <br>
Then build the main application<br> ``go build -o wb_app cmd/sub/main.go``<br>
Launch docker compose<br>``docker compose up -d`` <br>
Run application and then run publishing script <br>``./wb_app && cd cmd/pub && go run main.go``


