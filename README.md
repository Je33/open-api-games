# Open Api Games

Platform for integration of online games

## Service description

API service for integration external games to platform

## Project structure

### Business layer:

- `/domain`: data models, used in business logic
- `/service`: business logic, separated by logical use cases

### Technical layers:

- `/config`: configuration of application
- `/repository`: database storages
- `/provider`: external service providers
- `/transport`: interfaces for interaction with application

## Development

Scripts and tools for development and debugging can be found in `/Makefile`

To start dev server with all dependencies you'll need installed and running Docker Desktop (https://www.docker.com/products/docker-desktop/) and run:
```shell
make dc
```
This will up docker-composer.yml file with mongodb server in replica set mode (to support transactions) and instance of api application with hot reload support by changing code on the local address: http://localhost:8080.

Application code also contains the seed service, which fill initial data to database for testing on the first start, example of requests to test application functionality:

To check balance of test user:
```shell
curl --location 'http://localhost:8080/open-api-games/v1/games-processor' \
--header 'Sign: 02cc2f56e321fd39b911be6683b84076' \
--header 'Content-Type: application/json' \
--data '{
    "api": "balance",
    "data": {
        "gameSessionId": "FIRST_SESSION_UID",
        "currency": "USD"
    }
}'
```

To debit test user:
```shell
curl --location 'http://localhost:8080/open-api-games/v1/games-processor' \
--header 'Sign: 02a50cf0b98cf3baafd53bd53fe6cd51' \
--header 'Content-Type: application/json' \
--data '{
    "api": "debit",
    "data": {
        "gameSessionId": "FIRST_SESSION_UID",
        "currency": "USD",
        "amount": 100,
        "betId": "round-123"
    }
}'
```

And to credit test user
```shell
curl --location 'http://localhost:8080/open-api-games/v1/games-processor' \
--header 'Sign: 46c0f0cd2184b619a15cb8bd9a6d1e0d' \
--header 'Content-Type: application/json' \
--data '{
    "api": "credit",
    "data": {
        "gameSessionId": "FIRST_SESSION_UID",
        "currency": "USD",
        "amount": 100,
        "betId": "round-123"
    }
}'
```

## Testing

All the business layer logic covered by tests and can be run with:
```shell
make test
```

To regenerate all mocks:
```shell
make gen
```