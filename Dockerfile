FROM golang:1.12 AS builder

COPY . /workspace
WORKDIR /workspace

ENV CGO_ENABLED 0

RUN go build -o bin/check src/check/main.go \
 && go build -o bin/in src/in/main.go \
 && go build -o bin/out src/out/main.go


FROM alpine AS certs

RUN apk update && apk add ca-certificates


FROM busybox

COPY --from=builder /workspace/bin /opt/resource
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
