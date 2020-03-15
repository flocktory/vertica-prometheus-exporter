FROM golang:1.14-alpine as builder
WORKDIR /build
COPY . /build
RUN go mod download && go build -v 

FROM alpine:3.11
COPY --from=builder /build/vertica-prometheus-exporter /usr/app/vertica-prometheus-exporter

ENTRYPOINT  ["/usr/app/vertica-prometheus-exporter"]