.PHONY: all test clean

all: test build alfred

build:
	mkdir -p ./bin
	go build -o ./bin/alfredkeychain

alfred:
	rm -rf ./assets
	mkdir -p ./assets
	zip -j assets/KeychainPwd.alfredworkflow alfred_workflow/* bin/alfredkeychain

test:
	go test

clean:
	rm -rf ./bin
	rm -rf ./assets/*
