# ImageElevator

ImageElevator is a service designed to pull tar images from an FTP server and push them to a Docker registry. This project facilitates the automated transfer and management of Docker images in environments where FTP servers are used to store Docker tarballs.

## Features

- Fetch Docker tar images from an FTP server.
- Upload the fetched images to a specified Docker registry.
- Easy configuration through environment variables or configuration files.
- Robust logging using Zerolog for easy debugging and monitoring.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/KJone1/ImageElevator.git
    cd ImageElevator
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    go build -o image-elevator
    ```

## Configuration

ImageElevator can be configured using a configuration file or environment variables. The primary configuration settings include:

- `FTP_HOST`: The hostname of the FTP server.
- `FTP_USER`: The FTP username.
- `FTP_PASS`: The FTP password.
- `DOCKER_REGISTRY`: The Docker registry URL where images will be pushed.
- `LOG_LEVEL`: The logging level (e.g., `info`, `debug`, `error`).

You can also customize other settings such as the file pattern to match tar files and the working directory for image extraction.

## Usage

After configuring, you can start the service by running:

```bash
./image-elevator
