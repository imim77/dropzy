build:
	go build -o bin/dropzy
run: build
	./bin/dropzy
test:
	go test ./... -v