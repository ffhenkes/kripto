include kserver.env
export $(shell sed 's/=.*//' kserver.env)

all: build test

build:
	$(PACKAGE)

run:
	export PHRASE=avocado
	./kserver
