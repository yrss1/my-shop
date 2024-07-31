FROM golang:1.22.5-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./

COPY . .

RUN go mod download

RUN go build -o user-service ./user
RUN go build -o product-service ./product
RUN go build -o order-service ./order
RUN go build -o payment-service ./payment
RUN go build -o api-gateway-service ./api-gateway

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /workspace/user-service ./user-service
COPY --from=builder /workspace/product-service ./product-service
COPY --from=builder /workspace/order-service ./order-service
COPY --from=builder /workspace/payment-service ./payment-service
COPY --from=builder /workspace/api-gateway-service ./api-gateway-service

CMD ["./api-gateway-service"]

