FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin ./cmd

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/web /app/web
CMD [ "/app/bin" ]

