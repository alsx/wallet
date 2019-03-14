# Wallet
An example of app to make payments.

## Installation

Use the [make](https://www.gnu.org/software/make/) and [docker-compose](https://docs.docker.com/compose/) to set up Wallet.

```bash
make env-up
```

to check
```bash
docker-compose ps
```

```
   Name                 Command              State                     Ports
-----------------------------------------------------------------------------------------------
wallet_app_1   ./wallet                        Up      80/tcp
wallet_app_2   ./wallet                        Up      80/tcp
wallet_app_3   ./wallet                        Up      80/tcp
wallet_lb_1    /traefik --api --docker         Up      0.0.0.0:80->80/tcp, 0.0.0.0:8080->8080/tcp
db             docker-entrypoint.sh postgres   Up      0.0.0.0:5432->5432/tcp
```

to stop
```bash
make env-down
```

## Usage

```bash
$ # Request a list of  accounts.
$ curl http://localhost/v1/accounts/ -s | jq

```
```json
{
  "accounts": [
    {
      "id": "bob123",
      "balance": 86,
      "currency": ""
    },
    {
      "id": "alice456",
      "balance": 14.01,
      "currency": ""
    }
  ]
}
```

```bash
# Make a payment.
curl -d '{"from_account":"bob123","amount": 7, "to_account": "alice456"}' -H "Content-Type: application/json" -X POST http://localhost/v1/payments/
```

```bash
# Get a list of payments.
curl http://localhost/v1/payments/ -s | jq
```

```json
{
  "payments": [
    {
      "account": "alice456",
      "amount": 50,
      "to_account": "bob123",
      "direction": "outgoing"
    },
    {
      "account": "bob123",
      "amount": 50,
      "from_account": "alice456",
      "direction": "incoming"
    }
  ]
}
```

## Dashboard

Traefik dashboard is available at http://localhost:8080/dashboard/.


## Tests

To run test please use `make test` or `make test-race`.

For performance test: `make perf-test`.
This require [siage](https://www.joedog.org/siege-home/).
```bash
brew install siege
```
Example:
```
$ make perf-test 
siege -i -c50 -t60S --content-type "application/json" -f urls.txt
** SIEGE 4.0.4
** Preparing 50 concurrent users for battle.
The server is now under siege...
HTTP/1.1 404     0.10 secs:      19 bytes ==> GET  /404/
HTTP/1.1 404     0.10 secs:      19 bytes ==> GET  /404/
HTTP/1.1 404     0.12 secs:      19 bytes ==> GET  /404/
HTTP/1.1 404     0.02 secs:      19 bytes ==> GET  /404/
HTTP/1.1 500     0.22 secs:     131 bytes ==> POST http://localhost/v1/payments/
HTTP/1.1 200     0.27 secs:       3 bytes ==> POST http://localhost/v1/payments/
HTTP/1.1 500     0.27 secs:      83 bytes ==> POST http://localhost/v1/payments/
HTTP/1.1 200     0.29 secs:       3 bytes ==> POST http://localhost/v1/payments/
HTTP/1.1 200     0.33 secs:   22391 bytes ==> GET  /v1/payments/
HTTP/1.1 200     0.33 secs:   22391 bytes ==> GET  /v1/payments/
HTTP/1.1 500     0.34 secs:      55 bytes ==> POST http://localhost/v1/payments/
...
Transactions:                    699 hits
Availability:                  39.94 %
Elapsed time:                  21.23 secs
Data transferred:               4.99 MB
Response time:                  1.51 secs
Transaction rate:              32.93 trans/sec
Throughput:                     0.24 MB/sec
Concurrency:                   49.70
Successful transactions:         457
Failed transactions:            1051
Longest transaction:            8.37
Shortest transaction:           0.00
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
