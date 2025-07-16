FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Run tests
RUN go test -v

# Build the application (if there was a main function)
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage for testing
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy test files
COPY --from=builder /app/*.go ./

# This is primarily for testing the BRRRR calculation logic
CMD ["echo", "BRRRR evaluation container ready. Run 'go test -v' to execute tests."]