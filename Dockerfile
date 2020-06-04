FROM golang:1.14 as base
WORKDIR /go/src/github.com/hunterlong/discord
ADD . /go/src/github.com/hunterlong/discord
RUN go mod download && \
    go build -o discord

RUN mv discord /usr/local/bin/discord

FROM qmcgaw/youtube-dl-alpine
USER root

COPY --from=base /go/src/github.com/hunterlong/discord/discord /usr/local/bin/discord

ENV YOUTUBE "empty"
ENV DISCORD "empty"

RUN discord
