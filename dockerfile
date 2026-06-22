FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o watchkeeper ./cmd/api-server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/watchkeeper .

ENTRYPOINT [ "./watchkeeper" ]