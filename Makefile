include kripto.env
export $(shell sed 's/=.*//' kripto.env)

all: build test

deps:
	go get -v -t ./...

test:
	go test -v ./...

build:
	cd cmd/kserver && make -e build

run:
	cd cmd/kserver && make -e run

wrap: build
	docker build -t $(IMAGE):$(TAG) .

docker-clean:
	docker rm -f $(NAME) | true

docker-run: docker-clean
	docker run --net=host -d -ti --name=$(NAME) -e KRIPTO_ADDRESS=$(KRIPTO_ADDRESS) --volume=$(HOST_VOLUME):/data/secrets $(IMAGE):$(TAG)
	docker logs -f $(NAME)
