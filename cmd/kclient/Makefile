include kclient.env
export $(shell sed 's/=.*//' kclient.env)

all: build test

build:
	$(PACKAGE)

run:
	./kripto
