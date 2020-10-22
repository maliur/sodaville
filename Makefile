run:
	go build -o sodaville cmd/bot/main.go && ./sodaville

test:
	go test ./...

lint:
	golangci-lint run ./...

.PHONY: compile
compile:
	GOOS=linux GOARCH=amd64 go build -o sodaville-linux_amd64 ./cmd/bot/main.go
	GOOS=windows GOARCH=amd64 go build -o sodaville-windows_amd64.exe ./cmd/bot/main.go
	GOOS=linux GOARCH=arm go build -o sodaville-linux_arm ./cmd/bot/main.go

.PHONY: clean
clean:
	rm -f sodaville-linux_amd64 sodaville-windows_amd64.exe sodaville-linux_arm
