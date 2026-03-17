# First stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# getting binary file
RUN CGO_ENABLED=0 GOOS=linux go build -o shipment-service ./cmd/main.go

# Second stage
FROM alpine:latest

WORKDIR /root/

# copying binary file 
COPY --from=builder /app/shipment-service .

EXPOSE 50051

CMD ["./shipment-service"]