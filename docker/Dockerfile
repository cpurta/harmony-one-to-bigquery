FROM alpine:3.13

COPY ["artifacts/hmy-bq-import", "/usr/local/bin"]

ENTRYPOINT ["/usr/local/bin/hmy-bq-import"]
CMD ["backfill"]
