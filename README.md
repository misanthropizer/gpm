# gpm

## Overview

GNU Private Messenger (gpm) is a secure and private messaging application for nerds.

## Features

- End-to-end encryption using OpenPGP
- Configurable with personal OpenPGP keys
- Simple CLI for sending and receiving encrypted messages

## Requirements

- Go
- Docker (for Docker-based installation)

## Installation

### Go Build Instructions (Linux)

1. Clone the repository:
```sh
git clone https://github.com/misanthropizer/gpm.git
```
2. Navigate into the cloned directory:
```sh
cd gpm
```
3. Build the application using Go:
```sh
go build
```
4. Run gpm:
```sh
./gpm
```

### Docker Instructions (Debian)

1. Pull the Debian Docker image:
```sh
docker pull debian
```
2. Clone the gpm repository:
```sh
git clone https://github.com/misanthropizer/gpm.git
```
3. Navigate into the cloned directory:
```sh
cd gpm
```
4. Build the Docker container:
```sh
docker build -t gpm .
```
5. Run the Docker container:
```sh
docker run -it gpm
```

### Clone and Build Docker Container (Linux)

1. Clone the gpm repository:
```sh 
git clone https://github.com/misanthropizer/gpm.git
```
2. Navigate into the cloned directory:
```sh 
cd gpm
```
3. Build the Docker container:
```sh 
docker build -t gpm .
```
4. Run the Docker container:
```sh 
docker run -it gpm
```

## Usage
After launching gpm, follow the on-screen instructions to configure your private and public OpenPGP keys. You can then send and receive encrypted messages.
## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your changes.
## License
gpm is released under the MIT License. See the LICENSE file for more details.
