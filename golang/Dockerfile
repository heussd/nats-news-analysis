FROM golang:1.24.2-bullseye AS base

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download


# https://github.com/drone/ca-certs/blob/master/Dockerfile
FROM alpine:3.6 AS alpine
RUN apk add -U --no-cache ca-certificates


FROM base AS devcontainer

RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/nats-io/natscli/nats@latest

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apt-get update && \
    apt-get install -y locales && \
    rm -rf /var/lib/apt/lists/*
RUN echo "LC_ALL=en_US.UTF-8" >> /etc/environment && \
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen && \
    echo "LANG=en_US.UTF-8" > /etc/locale.conf && \
    locale-gen en_US.UTF-8


FROM base AS builder

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

ARG MAIN="keyword-matcher"
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin /app/cmd/${MAIN}


FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/goapp.bin /goapp.bin

CMD [ "./goapp.bin" ]
