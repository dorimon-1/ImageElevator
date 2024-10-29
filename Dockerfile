# Use the official Golang image to build the application
FROM golang:1.22.3 AS builder

# Install gpgme and btrfs dependencies
RUN apt-get update && apt-get install -y \
  libgpgme-dev \
  libbtrfs-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o image-elevator

# Start a new stage from scratch
FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
  libgpgme-dev \
  libbtrfs-dev
# Add ca-certificates
# RUN apk --no-cache add ca-certificates
# RUN apk --no-cache add libgpgme btrfs-progs

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/image-elevator .

COPY --from=builder /app/.env .

# Command to run the executable
CMD [ "./image-elevator" ]

