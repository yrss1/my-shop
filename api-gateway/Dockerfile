FROM golang:1.22.5-alpine as builder

WORKDIR /build

COPY . /build

RUN go mod download

COPY . .

RUN go build -o api-gateway-service .

FROM alpine:3.18 as hoster
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/api-gateway-service ./api-gateway-service

ENTRYPOINT ["./api-gateway-service"]