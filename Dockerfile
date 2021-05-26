FROM golang:1.16 AS builder

WORKDIR /go/src/fee_collector
RUN mkdir -p /build
COPY . .
RUN go build -race -ldflags "-extldflags '-static'" -o /build/fee_collector

FROM ubuntu:20.04 as fee_collector
RUN apt update && apt install ca-certificates -y && mkdir ./record

WORKDIR /btc_fee_collector

COPY --from=builder /go/src/fee_collector/config.yaml .
COPY --from=builder /build/fee_collector .

ENTRYPOINT ["/btc_fee_collector/fee_collector"]