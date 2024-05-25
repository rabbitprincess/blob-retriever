.PHONY=all build run

all:
	make build && make run

build:
	go build -o blob-retriever ./cmd/main.go

run:
	./blob-retriever