FROM golang:latest

WORKDIR /workdir

COPY ./main.go ./

COPY ./vendor ./vendor

COPY ./discord ./discord

COPY ./textToSpeech ./textToSpeech

COPY ./go.mod ./go.mod

COPY ./go.sum ./go.sum

RUN go build -o /app

CMD ["/app"]