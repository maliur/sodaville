run:
	go build -o sodaville cmd/main.go && ./sodaville

test:
	go test ./...
