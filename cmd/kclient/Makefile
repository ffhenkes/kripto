include kclient.env
export $(shell sed 's/=.*//' kclient.env)

PHRASE := $(if $(PHRASE),$(PHRASE),$(DEFAULT))

all: build test

build:
	$(PACKAGE)

run:
	./kripto
