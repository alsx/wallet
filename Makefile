FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

init:
	export GO111MODULE=on
	go mod vendor

test: init                           ## Run tests.
	go test -v ./...

test-race: init                      ## Run tests with race detector.
	go test -v -race ./...

perf-test:
	siege -i -c50 -t60S --content-type "application/json" -f urls.txt

format:                         ## Format source code.
	gofmt -w -s $(FILES_WTEST)
	goimports -local github.com/alsx/wallet -l -w $(FILES_WTEST)

env-down:
	docker-compose down --volumes --remove-orphans
	rm -rf logs

env-up:
	docker-compose up -d --build --scale  app=3

db-client:                      ## Connect to DB.
	docker exec -ti --user postgres db psql
