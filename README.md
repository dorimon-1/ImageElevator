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
    task setup
    ```

3. Build the project:
    ```bash
    task build
    ```
- `Run Tests`
  ```bash 
  task test
  ```

## Configuration

ImageElevator can be configured using a configuration file or environment variables. The primary configuration settings include:

- `REGISTRY`: The registry you wish to upload docker images to, for example: docker.io
- `REPOSITORY`: The repository you wish to upload to the docker images inside the registry, for example grafana
- `REPO_USERNAME`: The Username you are using to login to the registry
- `REPO_PASSWORD`: The Password you are using to login to the registry
- `DOCKER_CERT_PATH`: Path to a your docker certificate - Very Optional
- `TAR_REGEX`: The regex that matches your images pattern, leaving this empty will try to upload every file in your ftp server.
- `REGISTRY_BEARER_TOKEN`: An alternate, easier and more secure way to authenticate to your registry instead of using username and password.
- `SYNC_REGISTRIES`: Registries you wish to sync the image with for example: "docker.io, openshift.co"
- `SYNC_REGISTRIES_BEARER_TOKEN`: Bearer token for each registry
- `SAMPLE_RATE_IN_MINUTES`: The rate of your samples.
- `ZIP_REGEX`: The regex that matches your zip pattern, leaving this empty will try to upload every file in your ftp server.
- `ZIP_DESTINATION_PATH`: The destination folder (should be a mount nas etc) the zip will be copied to, leaving this empty will disable zip elevator.
- `FTP_HOST`: The hostname of the FTP server.
- `FTP_SERVER_PATH`: Working path inside the ftp server, default: "/"
- `FTP_USER`: The FTP username.
- `FTP_PASS`: The FTP password.
- `FTP_LOGGER_ENABLED`: Anything but "" to enable ftp logging
- `LOG_LEVEL`: The logging level (e.g., `info`, `debug`, `error`).

You can also customize other settings such as the file pattern to match tar files and the working directory for image extraction.

## Usage

After configuring, you can start the service by running:

```bash
./image-elevator
```
