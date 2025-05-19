FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-app ./cmd

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/chat-app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./chat-app"]
