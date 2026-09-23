# Build stage
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate


# Runtime stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/migrate ./migrate

EXPOSE 8080

CMD ["./server"]