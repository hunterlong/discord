FROM golang:1.14.2-alpine as base
RUN apk add --no-cache libstdc++ gcc g++ make git ca-certificates linux-headers wget curl jq

WORKDIR /go/src/github.com/hunterlong/discord
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -o discord . && \
    chmod +x discord

FROM qmcgaw/youtube-dl-alpine
USER root

RUN apk --no-cache add curl jq ca-certificates linux-headers bash

WORKDIR /root/
COPY --from=base /go/src/github.com/hunterlong/discord/discord .

ENV YOUTUBE "empty"
ENV DISCORD "empty"
ENV CHANNEL_ID "empty"
ENV GUILD_ID "empty"
ENV CHANNELS "empty"
ENV LIMIT 25
ENV BITRATE 64

ENTRYPOINT /root/discord
