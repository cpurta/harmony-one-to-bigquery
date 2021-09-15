FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/cpurta/harmony-one-to-bigquery
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./
COPY ./go.sum ./
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/hmy-bq-import ./cmd/hmy-bq-import

FROM alpine:3.13
RUN mkdir -p /etc/hmy/creds
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder ["/go/src/github.com/cpurta/harmony-one-to-bigquery/bin/hmy-bq-import", "/usr/local/bin"]
COPY ["docker/artifacts/harmonyone-gcp-bigquery.json", "/etc/hmy"]
ENTRYPOINT ["/usr/local/bin/hmy-bq-import"]
CMD ["backfill"]
