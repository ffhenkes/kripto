include kripto.env
export $(shell sed 's/=.*//' kripto.env)

all: build test

build:
	$(PACKAGE)

run:
	./kserver
