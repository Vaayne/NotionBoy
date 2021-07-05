
build:
	docker build -t ghcr.io/vaayne/notion-boy .

init:
	if ! which pre-commit > /dev/null; then sudo pip install pre-commit; fi
	pre-commit install

static: init
	pre-commit run --all-files

run:
	go run ./main.go

rund: build
	docker run --rm --env-file=.env ghcr.io/vaayne/notion-boy

push:
	GOOS=linux go build -o ./app .
	scp -P 121  /Users/vaayne/Github/Notion-Boy/app   ubuntu@121.36.78.252:/home/ubuntu
