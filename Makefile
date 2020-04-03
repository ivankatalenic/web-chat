build:
	go build -v .

make lint:
	golint ./...

make test:
	go test -v -count=1 ./...

check:
	golint ./...
	go test -v -count=1 ./...
