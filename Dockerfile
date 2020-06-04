FROM golang:1.14 as base
RUN apt install -y libstdc++ gcc g++ make git ca-certificates linux-headers wget curl jq

WORKDIR /go/src/github.com/hunterlong/discord
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -o discord . && \
    chmod +x discord

FROM qmcgaw/youtube-dl-alpine
USER root

RUN apk --no-cache add curl jq ca-certificates linux-headers

WORKDIR /root/
COPY --from=base /go/src/github.com/hunterlong/discord/discord .

ENV YOUTUBE "empty"
ENV DISCORD "empty"
ENV CHANNEL_ID "empty"
ENV GUILD_ID "empty"
ENV CHANNELS "empty"

ENTRYPOINT /bin/sh
CMD ["/root/discord"]
