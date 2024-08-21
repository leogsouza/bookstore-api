# Official Go Apline Base Image
FROM golang:1.22-alpine as builder

# Create The Application Directory
WORKDIR /app

# Copy and Download Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy The Application Source & Build
COPY . .

ENV CGO_ENABLED=0
RUN go build -o main .

# Final Image Creation Stage
FROM alpine:3.19

WORKDIR /root/

# Copy The Built Binary
COPY --from=builder /app/main .

# Expose the port
CMD ["./main"]