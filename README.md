# URL Shortener

A fast and efficient URL shortening service built in Go, featuring persistent storage with SQLite and in-memory caching.

## Features
- URL shortening using CRC32 hash
- SQLite database tracking of shortened URLs
- Request origin tracking
- In-memory cache for fast URL lookups
- Clean web interface with dark theme
- HTMX for dynamic updates

## Getting Started

### Prerequisites
- Go 1.19 or higher
- SQLite

### Installation
```bash
git clone <repository-url>
cd url-shortener
```

### Running the Service
```bash
go run .
```
The server will start on `localhost:8080`

### Running Tests
```bash
go test -v
```

## Usage
- Visit the web interface at `http://localhost:8080`
- Submit a URL to receive a shortened version
- Access shortened URLs via `/q/<short-code>`
- View recent URLs and their statistics

## Database
The SQLite database (`urls.sql`) tracks:
- Original URL
- Shortened code
- Creation timestamp
- Requester IP address

## API Endpoints
- `POST /s` - Create short URL
- `GET /u` - List all URLs
- `GET /q/<short-code>` - Redirect to original URL

## Tech Stack
- Go
- SQLite
- HTMX
- CSS

Happy URL shortening! 🚀
