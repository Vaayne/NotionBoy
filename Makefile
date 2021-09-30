.PHONY: build clean

build: build-wxgzh build-netlify

init:
	if ! which pre-commit > /dev/null; then sudo pip install pre-commit; fi
	pre-commit install

static: init
	pre-commit run --all-files

run:
	go run ./main.go

rund: build
	docker run --rm ghcr.io/vaayne/notion-boy

build-netlify:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/notion cmd/notion/main.go

build-wxgzh:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/wxgzh cmd/wxgzh/main.go

build-docker:
	docker build -t ghcr.io/vaayne/notion-boy .

clean:
	rm -rf ./bin ./vendor go.sum
