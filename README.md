# caching-proxy
A small caching proxy to cache http/https requests

## Installation

### Dev

```bash
go build
go run . -p 3000 -origin https://example.com
```

### User

```bash
go install github.com/Fidnix/caching-proxy@v0.1.1
```

> [!NOTE]
> Maybe tou have to check your GOBIN env variable with `go env`

## Usage

```bash
caching-proxy --port 3000 -origin https://example.com
```

When execute the program, it creates a server to listen `https://example.com` requests behind `localhost:3000`

Each response returned will have a **x-cache header** with values: **HIT**, **MISS**; depending on if it's cached

## TODO

### Funcionality

- [ ] Cache content by query params
- [ ] Cache content by request body
- [ ] Allow log option
- [ ] Verify args

### Project

- [ ] Modularize project
- [ ] Tests