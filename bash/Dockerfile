# versions > 0.2.2 causes bash switch statement to not work anymore?!
ARG NATSCLI_VERSION=0.2.2

FROM alpine:latest AS build
ARG NATSCLI_VERSION
RUN apk add --no-cache --update wget sed
WORKDIR /download
RUN arch=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/) && \
     wget https://github.com/nats-io/natscli/releases/download/v${NATSCLI_VERSION}/nats-${NATSCLI_VERSION}-linux-$arch.zip \
     --progress=dot:giga \
     --output-document 'natscli.zip' && \
     unzip natscli.zip && \
     mv nats-*/* ..


FROM alpine:latest
RUN apk add --no-cache --update bash

COPY --from=build /nats /
COPY --chmod=755 feed.sh feed-urls.json /

ENTRYPOINT [ "/feed.sh" ]