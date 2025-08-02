FROM golang:1.23.9

RUN go install github.com/nats-io/natscli/nats@v0.2.4

COPY feed.sh /

ENTRYPOINT [ "/feed.sh" ]