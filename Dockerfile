FROM golang:1.14 as base

WORKDIR /go/src/github.com/hunterlong/discord
ADD go.mod /go/src/github.com/hunterlong/discord
ADD go.sum /go/src/github.com/hunterlong/discord
RUN go mod download

ADD . /go/src/github.com/hunterlong/discord
RUN go build -o discord && \
    chmod +x discord

FROM qmcgaw/youtube-dl-alpine
USER root

COPY --from=base /go/src/github.com/hunterlong/discord/discord /usr/local/bin/discord

ENV YOUTUBE "empty"
ENV DISCORD "empty"

CMD /usr/local/bin/discord
