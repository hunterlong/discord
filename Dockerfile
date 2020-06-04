FROM golang:1.14 as base

WORKDIR /go/src/github.com/hunterlong/discord
ADD go.mod /go/src/github.com/hunterlong/discord
ADD go.sum /go/src/github.com/hunterlong/discord
RUN go mod download

ADD . .
RUN go build -o discord . && \
    chmod +x discord

FROM qmcgaw/youtube-dl-alpine
USER root

WORKDIR /root/
COPY --from=base /go/src/github.com/hunterlong/discord/discord .

ENV YOUTUBE "empty"
ENV DISCORD "empty"
ENV CHANNEL_ID "empty"
ENV GUILD_ID "empty"
ENV CHANNELS "empty"

ENTRYPOINT "/root/discord"
