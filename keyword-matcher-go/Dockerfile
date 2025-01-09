FROM golang:1.22-bullseye as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin ./cmd


# https://github.com/drone/ca-certs/blob/master/Dockerfile
FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates


FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/goapp.bin /goapp.bin

CMD [ "./goapp.bin" ]
