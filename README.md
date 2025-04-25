# URL Shortener API

A simple URL shortener API built with Go, using bbolt for storage and apiculi for API routing.

## Features

- Shorten URLs with optional custom pretty names
- Redirect from shortened URLs to original URLs
- API key protection for URL creation
- Persistent storage with bbolt
- OpenAPI documentation
- Configurable via environment variables

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| HOST | Host to bind the server to | localhost |
| PORT | Port to listen on | 8080 |
| DB_PATH | Path to the bbolt database file | shortener.db |
| API_KEY | API key for protected endpoints | your-api-key-here |
| BASE_URL | Base URL for shortened links | http://localhost:8080 |

### Running the Application

```bash
go run main.go
```

Or using the provided npm script:

```bash
npm start
```

### Building the Application

```bash
go build -o url-shortener
```

## API Endpoints

### Redirect to Original URL

```
GET /{code}
GET /{code}/{pretty-name}
```

### Create a Shortened URL

```
POST /api/urls
```

Headers:
```
X-API-Key: your-api-key-here
Content-Type: application/json
```

Request Body:
```json
{
  "original_url": "https://example.com/some/very/long/url",
  "pretty_name": "example"  // Optional
}
```

Response:
```json
{
  "code": "a1b2c3d4",
  "short_url": "http://localhost:8080/a1b2c3d4/example",
  "original_url": "https://example.com/some/very/long/url",
  "pretty_name": "example"
}
```

### List All URLs

```
GET /api/urls
```

Headers:
```
X-API-Key: your-api-key-here
```

## License

This project is licensed under the MIT License.