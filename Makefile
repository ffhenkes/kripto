include docker.env
export $(shell sed 's/=.*//' docker.env)

all: build test

deps:
	go get -v -t ./...

test: deps
	go test -v -cover ./...

build:
	cd cmd/kserver && make -e build

build-cli:
	cd cmd/kclient && make -e build

run: build
	cd cmd/kserver && make -e run

cli: build-cli
	cd cmd/kclient && make -e run

certificates:
	./scripts/self-sign.sh

signature:
	./scripts/signature.sh

wrap: build build-cli
	docker build -t $(IMAGE):$(TAG) .

docker-clean:
	docker rm -f $(NAME) | true

docker-run: docker-clean
	docker run --net=host -d -ti --name=$(NAME) --env-file=$(NAME).env --volume=$(SECRETS_VOLUME):/secrets --volume=$(AUTH_VOLUME):/authdb $(IMAGE):$(TAG)
	docker logs -f $(NAME)
