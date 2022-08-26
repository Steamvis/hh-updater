FROM golang:1.18 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

FROM ubuntu:xenial

RUN apt-get update && apt-get install -y chromium-browser

RUN apt-get clean autoclean && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}

WORKDIR /
COPY --from=build /usr/local/bin/app /
CMD ["/app"]