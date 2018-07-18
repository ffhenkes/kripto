all: build test

deps:
	go get -v -t ./...

test:
	go test -v ./...

build:
	cd cmd/kserver && make -e build

run:
	cd cmd/kserver && make -e run
