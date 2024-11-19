# Build stage
FROM golang:1.20-alpine AS build

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o /ordering-system ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary from the build stage
COPY --from=build /ordering-system .


# Expose the application port
EXPOSE 8080

# Run the Go app
CMD [ "db:5432", "--", "./ordering-system"]
