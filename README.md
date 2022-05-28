# Distributed Tracing Golang Sample

This project demonstrates how to instrument distrubted tracing in golang application with the help of below microservices:

- order-service
- payment-service
- user-service

For this you would need the following:

- Go (version >= 1.16): For installation see [getting started](https://go.dev/doc/install)
- MySQL 8: Download the MySQL community version from [here](https://dev.mysql.com/downloads/mysql/)
- `serve` for the frontend. For installation see: [https://www.npmjs.com/package/serve](https://www.npmjs.com/package/serve)
- [Signoz](https://signoz.io/)

## Tracing flow
![distributed_tracing_app_otel_signoz](https://user-images.githubusercontent.com/83692067/170809287-f930245e-55b9-4646-8e71-74dfe74e036e.png)


## Running the code

Start the signoz server following the instructions:

```sh
git clone -b main https://github.com/SigNoz/signoz.git
cd signoz/deploy/
./install.sh
```

Configuration for microservices can be updated in .env file

```
# service config
USER_URL=localhost:8080
PAYMENT_URL=localhost:8081
ORDER_URL=localhost:8082

# database config
SQL_USER=root
SQL_PASSWORD=password
SQL_HOST=localhost:3306
SQL_DB=signoz

# telemetry config
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
INSECURE_MODE=true
```

Start individual microservices using below commands

1. User Service

```sh
go run ./users
```

2. Payment Service

```sh
go run ./payment
```

3. Order Service

```sh
go run ./order
```

Start the frontend using following command. For installation of `serve` see: [https://www.npmjs.com/package/serve](https://www.npmjs.com/package/serve)

```sh
serve -l 5000 frontend
```

View traces and metrics at http://localhost:3301/
