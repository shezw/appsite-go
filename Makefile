.PHONY: build run test clean coverage

APP_NAME = appsite-monolith
CMD_PATH = ./cmd/$(APP_NAME)
CONFIG_PATH = ./configs/config.yaml

build:
	go build -v -o bin/$(APP_NAME) $(CMD_PATH)

run: build
	./bin/$(APP_NAME)

test:
	go test -v ./...

coverage:
	./coverage.sh all

clean:
	rm -rf bin/ coverage.out

test-api: build
	@nohup ./bin/$(APP_NAME) > server.log 2>&1 & echo $$! | ./tests/scripts/test_apis.sh
