# build executable
FROM golang:1.23 AS builder

WORKDIR /build

# Copy and download dependencies using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the applications
RUN go build -o /app/fdc-client main/main.go

FROM debian:latest AS execution

WORKDIR /app
COPY --from=builder /app/fdc-client .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


CMD ["./fdc-client" ]
