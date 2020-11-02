FROM golang:1.13-alpine AS build
WORKDIR /tsm
COPY . .
RUN go build .

FROM alpine:latest
RUN addgroup -g 10000 -S tsm && \
    adduser  -u 10000 -S tsm -G tsm -H -s /bin/false && \
    apk --no-cache add su-exec
WORKDIR /tsm
COPY --from=build --chown=tsm:tsm /tsm/tsm /tsm/bin/start.sh /tsm/
CMD ["su-exec", "tsm", "sh", "/tsm/start.sh"]
EXPOSE 8080/tcp
