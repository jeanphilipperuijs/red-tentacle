# Red Tentacle

Go Reverse Proxy with Dynamic Backends

This is a simple reverse proxy written in Go that forwards incoming requests to multiple backend servers. The proxy sends the request to all configured backends and returns the response from the first successful backend. The list of backend servers can be updated dynamically at runtime via an API endpoint.

## Features

- Forwards incoming requests to multiple backend servers simultaneously.
- Returns the response from the first backend that responds successfully.
- Dynamically update the list of backend servers without restarting the service.
- Configurable backend servers via environment variables.

## Setup and Usage
1. Set up backend servers

Export the BACKEND_SERVERS environment variable with a comma-separated list of backend URLs:

```bash
export BACKEND_SERVERS="http://backend1.example.com,http://backend2.example.com"
```

2. Run the proxy

```bash
go run proxy.go
```
This will start the proxy server on localhost:8080.

3. Update backends at runtime

You can dynamically update the list of backends by sending a POST request to /update-backends:

```bash
curl -X POST "http://localhost:8080/update-backends?backends=http://new-backend1.com,http://new-backend2.com"
```

4. Forward requests

Send requests to the proxy, and it will forward them to the configured backend servers, returning the first successful response:

```bash
curl http://localhost:8080/some-endpoint
```

## License

This project is licensed under the MIT License.