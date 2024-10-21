# Red Tentacle

(because purple tentacle wants to take over the world)

This is a simple reverse proxy written in Go that forwards incoming requests to multiple backend servers. The proxy sends requests to all backends concurrently and returns the response of the first successful backend (HTTP status in the 2xx range). Backend servers are configured through an environment variable.

## Features

- Sends requests to multiple backends concurrently.
- Returns the response of the first successful backend.
- Supports dynamic configuration of backends through an environment variable.
- Logs errors and response statuses from the backends.
- Configurable request timeout for each backend.
- Update backend servers at runtime

## Environment Variables

- `BACKEND_SERVERS`: A comma-separated list of backend URLs to which the proxy will forward requests.

The advantage of this that you can easily change the backends without changing the code (or lua script)

### Example:

```bash
export BACKEND_SERVERS="http://backend1.example.com,http://backend2.example.com,http://backend3.example.com"
```

or updating at runtime 
```bash
curl -X POST http://localhost:8081/-update-backends\?backends\=https://dc1.vptech.eu,https://dc2.vptech.eu,https://dc3.vptech.eu
```

## Installation

### Clone the repository:

```bash
git clone https://github.com/yourusername/golang-reverse-proxy.git
cd golang-reverse-proxy
```

### Set up Go:

Make sure you have Go installed. You can download it from here.

Set the environment variable:

```bash
export BACKEND_SERVERS="http://backend1.example.com,http://backend2.example.com,http://backend3.example.com"
```

### Run the proxy:

```bash
go run proxy.go
```

## Usage

After starting the server, you can send HTTP requests to the proxy (listening on port 8080 by default), and the proxy will forward the request to all the backends defined in the BACKEND_SERVERS environment variable.

```bash
curl http://localhost:8080/api/some-endpoint
```

The proxy will forward the request to all configured backends and return the response from the first backend that responds successfully.

## Configuration

You can modify the timeout for backend requests by changing the Timeout value in the http.Client configuration inside the code.

## Error Handling

If all backend requests fail (no backend returns a 2xx status code), the proxy will return a 502 Bad Gateway response to the client.
Logs are printed to the console for each request, including errors and responses from each backend.




## Alternatives

Depending on your needs there is another more or less similar approach.

### nginx + lua

Using nginx with lua
```bash
sudo apt-get install -y openresty
```

```lua
http {
    server {
        listen 80;
        server_name yourproxy.com;

        location / {
            content_by_lua_block {
                local http = require "resty.http"
                local backend_urls = {
                    "http://backend1.example.com",
                    "http://backend2.example.com",
                    "http://backend3.example.com"
                }

                -- Send requests to all backends
                for _, backend in ipairs(backend_urls) do
                    local httpc = http.new()
                    httpc:request_uri(backend .. ngx.var.request_uri, {
                        method = ngx.req.get_method(),
                        body = ngx.req.get_body_data(),
                        headers = ngx.req.get_headers(),
                    })
                end

                -- Return a response
                ngx.say("Request forwarded to all backends")
            }
        }
    }
}
```

### Nginx with `sub_filter` for WebSocket or Long Polling

If you're dealing with WebSocket connections or asynchronous communication, you could use Nginx's sub_filter module or similar techniques to propagate responses from multiple servers back to the client. However, this is more complex and may require custom logic based on the application type.

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
