FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go build -v -o adservers .

FROM alpine:3.20

COPY --from=builder /app/adservers /usr/local/bin/adservers

WORKDIR /app

CMD ["adservers"]
