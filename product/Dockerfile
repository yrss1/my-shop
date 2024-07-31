FROM golang:1.22.5-alpine as builder

WORKDIR /build

COPY . /build

RUN go mod download

COPY . .

RUN go build -o product-service .

FROM alpine:3.18 as hoster
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/product-service ./product-service
COPY --from=builder /build/db/migrations ./db/migrations

ENTRYPOINT ["./product-service"]