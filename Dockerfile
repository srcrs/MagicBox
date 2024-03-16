FROM golang:1.20.4-alpine3.18 AS builder

COPY . /src
WORKDIR /src

RUN go build -ldflags "-s -w" -o ./bin/MagicBox .

FROM srcrs/real-chrome-stable:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

ENTRYPOINT ["./MagicBox"]
