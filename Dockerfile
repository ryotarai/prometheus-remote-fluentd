FROM golang:1.11 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /usr/bin/prometheus-remote-fluentd .

###############################################

FROM ubuntu:16.04
COPY --from=builder /usr/bin/prometheus-remote-fluentd /usr/bin/prometheus-remote-fluentd
ENTRYPOINT ["/usr/bin/prometheus-remote-fluentd"]
