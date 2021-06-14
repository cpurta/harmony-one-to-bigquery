FROM alpine:3.13

RUN mkdir -p /etc/hmy/creds

COPY ["artifacts/hmy-bq-import", "/usr/local/bin"]

COPY ["artifacts/harmonyone-gcp-bigquery.json", "/etc/hmy"]

ENTRYPOINT ["/usr/local/bin/hmy-bq-import"]
CMD ["backfill"]
