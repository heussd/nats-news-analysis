FROM golang:latest

RUN go install github.com/nats-io/natscli/nats@latest

COPY feed.sh /

ENTRYPOINT [ "/feed.sh" ]