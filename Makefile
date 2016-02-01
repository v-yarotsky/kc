.PHONY: all test clean

all:
	mkdir -p ./bin
	go build -o ./bin/alfredkeychain

test:
	go test

clean:
	rm -rf ./bin
