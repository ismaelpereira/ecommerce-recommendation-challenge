.PHONY: build run test race cover fmt tidy clean

APP_NAME := ecommerce-recommendation
MAIN := .

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) .

run:
	go run $(MAIN)

test:
	go test -v ./...

test-integration:
	@if [ -z "$$BIGTABLE_EMULATOR_HOST" ]; then \
		echo "BIGTABLE_EMULATOR_HOST is not set"; \
		exit 1; \
	fi
	go test -v -tags=integration ./...

race:
	go test -race ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -rf bin coverage.out