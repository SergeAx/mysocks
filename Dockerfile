FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN mkdir /app
COPY proxy /app/proxy

RUN chmod +x /app/proxy

ENTRYPOINT ["/app/proxy"]