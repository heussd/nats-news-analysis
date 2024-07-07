FROM golang:1.22-bullseye as base

ENV GO111MODULE=on


FROM base as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd


FROM builder as keyword-matcher-builder
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin ./cmd/keyword_matcher


FROM builder as raindrop-integration-builder
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin ./cmd/raindrop-integration


FROM builder as rss-article-url-feeder-builder
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -o goapp.bin ./cmd/rss-article-url-feeder



FROM scratch as raindrop-integration
COPY --from=raindrop-integration-builder /app/goapp.bin /goapp.bin
CMD [ "./goapp.bin" ]


FROM scratch as rss-article-url-feeder
COPY --from=rss-article-url-feeder-builder /app/goapp.bin /goapp.bin
CMD [ "./goapp.bin" ]


FROM scratch as keyword-matcher
COPY --from=keyword-matcher-builder /app/goapp.bin /goapp.bin
CMD [ "./goapp.bin" ]

