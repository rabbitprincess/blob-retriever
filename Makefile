.PHONY=all build run

all:
	make build && make run

build:
	go build -o blob-retriever ./cmd/

run:
	./blob-retriever