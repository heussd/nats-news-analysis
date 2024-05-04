FROM golang:1.22

RUN go install github.com/nats-io/natscli/nats@latest

COPY feed.sh /

ENTRYPOINT [ "/feed.sh" ]