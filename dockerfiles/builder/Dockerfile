# Get an image that can build the project at all time
FROM golang:1.22.3

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

ENTRYPOINT [ "/bin/bash" ]
