# ====================
# Build stage
# ====================
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ====================
# Run stage
# ====================
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml .

EXPOSE 8080
ENV GIN_MODE=release

ENTRYPOINT ["/app/main"]
