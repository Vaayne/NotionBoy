FROM golang:1.18 as build_base

WORKDIR /go/src/github.com/Vaayne/Notion-Boy
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o ./app cmd/notionboy/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai
WORKDIR /service/

RUN touch data.db

COPY --from=build_base /go/src/github.com/Vaayne/Notion-Boy/app ./app
COPY settings.yaml settings.yaml

CMD ["./app"]
