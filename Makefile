run:
	go build -o sodaville cmd/bot/main.go && ./sodaville

test:
	go test ./...

lint:
	golangci-lint run ./...
