
# ZipServer-Go

ZipServer-Go is a lightweight Golang web server that serves static files directly from a ZIP file without extracting it. The server optionally supports an SSL/TLS layer for secure communication.

## Features

- Serve static files directly from a ZIP file without extraction.
- Optional SSL/TLS support with client certificate validation.
- Configurable IP address and port for the server.

## Installation

To install the dependencies, you'll need to have Go installed on your machine. Then, run:

```sh
go get github.com/octopwn/zipserver-go
```

## Usage

To run the server, use the following command:

```sh
go run main.go [options] <path_to_zipfile>
```

### Options

- `-a <address>`: IP/hostname to listen on. Default is `127.0.0.1`.
- `-p <port>`: Port to listen on. Default is `8000`.
- `--ssl-cert <file>`: Path to the SSL certificate file.
- `--ssl-key <file>`: Path to the SSL key file.
- `--ssl-ca <file>`: Path to the CA certificate file for client certificate validation.

### Examples

1. **Basic usage without TLS:**

    ```sh
    go run main.go -a 127.0.0.1 -p 8080 /path/to/static_files.zip
    ```

2. **Usage with TLS:**

    ```sh
    go run main.go -a 127.0.0.1 -p 8443 --ssl-cert /path/to/cert.pem --ssl-key /path/to/key.pem /path/to/static_files.zip
    ```

3. **Usage with TLS and client certificate validation:**

    ```sh
    go run main.go -a 127.0.0.1 -p 8443 --ssl-cert /path/to/cert.pem --ssl-key /path/to/key.pem --ssl-ca /path/to/ca.pem /path/to/static_files.zip
    ```

## Example Output

If you run the server with the following command:

```sh
go run main.go -p 8000 /path/to/static_files.zip
```

You can access the contents of the ZIP file in your browser at `http://localhost:8000/`.

