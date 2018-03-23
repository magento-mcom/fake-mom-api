FROM golang:1.10-alpine

WORKDIR /go/src/app
COPY . .

ENV GIT_TERMINAL_PROMPT 1

RUN apk add --no-cache git go \
    && go get -d -v ./... \
    && apk del git

RUN go install -v ./...

CMD ["app"]

EXPOSE 24213