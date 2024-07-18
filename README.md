# gpm

## Overview

GNU Private Messenger (gpm) is a secure and private messaging application for nerds.

### Go Build Instructions (Linux)

1. Clone the repository:
   git clone https://github.com/misanthropizer/gpm.git

2. Navigate into the cloned directory:
   cd gpm

3. Build the application using Go:
   go build

4. Run gpm:
   ./gpm

### Docker Instructions (Debian)

1. Pull the Debian Docker image:
   docker pull debian

2. Clone the gpm repository:
   git clone https://github.com/misanthropizer/gpm.git

3. Navigate into the cloned directory:
   cd gpm

4. Build the Docker container:
   docker build -t gpm .

5. Run the Docker container:
   docker run -it gpm

### Clone and Build Docker Container (Linux)

1. Clone the gpm repository:
   git clone https://github.com/misanthropizer/gpm.git

2. Navigate into the cloned directory:
   cd gpm

3. Build the Docker container:
   docker build -t gpm .

4. Run the Docker container:
   docker run -it gpm
