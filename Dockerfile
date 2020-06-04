FROM qmcgaw/youtube-dl-alpine
USER root
RUN apk add --no-cache git make musl-dev go linux-headers
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR ${GOPATH}/src/github.com/hunterlong/discord
ADD . .

RUN go mod download \
    go build -o discord

RUN mv discord /usr/local/bin/discord

ENV YOUTUBE "empty"

RUN discord
