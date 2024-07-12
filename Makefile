.PHONY=all build run

all:
	make build && make run

build:
	CGO_ENABLED=0 go build -o blob-retriever ./cmd/

run:
	cp -n .env ./cmd/ && nohup ./blob-retriever > output.log 2>&1 &