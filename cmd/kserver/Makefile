include kserver.env
export $(shell sed 's/=.*//' kserver.env)

PHRASE := $(if $(PHRASE),$(PHRASE),$(DEFAULT))

all: build test

build:
	$(PACKAGE)

run:
	export PHRASE=avocado
	./kserver
