.PHONY: build clean

IMAGE_NAME ?= ghcr.io/vaayne/notionboy
IMAGE_TAG ?= latest

build: build-wxgzh build-netlify

init:
	if ! which pre-commit > /dev/null; then pip install pre-commit; fi
	pre-commit install

static: init
	pre-commit run --all-files

run:
	go run ./cmd/notionboy/main.go

rund: build-docker
	docker run --rm ${IMAGE_NAME}:${IMAGE_TAG}

build-docker:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

clean:
	rm -rf ./bin ./vendor go.sum
