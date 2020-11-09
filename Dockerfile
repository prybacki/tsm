FROM golang:alpine AS build
WORKDIR /tsm
COPY . .
RUN go build .

FROM alpine:latest
USER 1000
COPY --from=build /tsm/tsm /tsm/
CMD ["/tsm/tsm"]