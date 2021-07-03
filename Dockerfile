FROM golang:1.16 as build_base

WORKDIR /go/src/github.com/Vaayne/Notion-Boy
ENV CGO_ENABLED=0
COPY . .
RUN go build -o ./app .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /service/
COPY --from=build_base /go/src/github.com/Vaayne/Notion-Boy/app ./app
CMD ["./app"]
