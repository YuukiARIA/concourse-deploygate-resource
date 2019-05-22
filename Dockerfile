FROM golang:1.12 AS builder

COPY . /workspace
WORKDIR /workspace

RUN go build -o bin/check src/check/main.go \
 && go build -o bin/in src/in/main.go \
 && go build -o bin/out src/out/main.go

FROM busybox

COPY --from=builder /workspace/bin /opt/resource
