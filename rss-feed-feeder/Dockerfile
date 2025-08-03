FROM golang:1.23.9

# versions > 0.2.2 causes bash switch statement to not work anymore?!
RUN go install github.com/nats-io/natscli/nats@v0.2.2

COPY feed.sh /

ENTRYPOINT [ "/feed.sh" ]