# Build stage
FROM golang:1.23-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 go build -o cassandra-example

# Run stage
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/cassandra-example /app/cassandra-example
CMD ["/app/cassandra-example"]
