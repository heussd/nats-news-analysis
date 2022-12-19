FROM golang:1.19-bullseye as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin ./cmd


FROM scratch
COPY --from=builder /app/goapp.bin /goapp.bin

CMD [ "./goapp.bin" ]
