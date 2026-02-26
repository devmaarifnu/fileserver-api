# CDN File Server API

A high-performance, tag-based file server API built with Go and Gin framework. This server provides secure file upload, download, list, and delete operations with token-based authentication and public/private file access control.

## Features

- ✅ **Token-based Authentication** - Multiple tokens with granular permissions
- ✅ **File Upload** - Tag-based file organization
- ✅ **Public/Private Files** - Fine-grained access control per file
- ✅ **File Download/View** - Direct file serving with caching
- ✅ **File List** - Filtering, searching, and pagination support
- ✅ **File Delete** - Secure file deletion with authorization
- ✅ **CORS Enabled** - Frontend-friendly configuration
- ✅ **Comprehensive Logging** - JSON/text format with rotation
- ✅ **Response Compression** - Gzip compression support
- ✅ **Clean Code** - Well-structured, maintainable codebase

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin Framework** - HTTP web framework
- **Viper** - Configuration management
- **Logrus** - Structured logging
- **Lumberjack** - Log rotation
- **UUID** - Unique filename generation

## Project Structure

```
cdn-fileserver/
├── cmd/
│   └── api/
│       └── main.go                # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loader
│   ├── models/
│   │   └── file_meta.go           # File metadata model
│   ├── handlers/
│   │   ├── upload.go              # Upload handler
│   │   ├── download.go            # Download handler
│   │   ├── list.go                # List handler
│   │   ├── delete.go              # Delete handler
│   │   └── health.go              # Health check handler
│   ├── services/
│   │   ├── file_service.go        # File operations business logic
│   │   └── storage_service.go     # Storage management
│   ├── middleware/
│   │   ├── auth.go                # Token authentication
│   │   ├── logger.go              # Request logging
│   │   ├── cors.go                # CORS configuration
│   │   └── recovery.go            # Panic recovery
│   ├── utils/
│   │   ├── response.go            # Response helpers
│   │   ├── validator.go           # Input validation
│   │   ├── filename.go            # Filename generator
│   │   └── file.go                # File utilities
│   └── routes/
│       └── routes.go              # API routes registration
├── pkg/
│   └── logger/
│       └── logger.go              # Logger configuration
├── storage/                       # File storage directory
├── logs/                          # Application logs
├── config.yaml                    # Configuration file
├── config.example.yaml            # Configuration example
├── Makefile                       # Build commands
└── README.md
```

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

1. **Clone the repository**
   ```bash
   cd cdn-fileserver
   ```

2. **Install dependencies**
   ```bash
   make install
   # or
   go mod download
   ```

3. **Configure the application**
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your settings
   ```

4. **Generate strong tokens** (for production)
   ```bash
   openssl rand -hex 32
   ```

5. **Run the application**
   ```bash
   # Development mode
   make run
   # or
   go run cmd/api/main.go
   ```

## Configuration

Edit `config.yaml` to configure the application:

```yaml
app:
  name: "CDN File Server"
  env: "development"  # development, staging, production
  port: 8080
  version: "1.0.0"
  domain: "cdn.maarifnu.or.id"

storage:
  base_path: "./storage"
  max_file_size: 52428800  # 50MB
  allowed_extensions:
    - jpg
    - jpeg
    - png
    - pdf
    # ... more extensions

tokens:
  - id: "token_001"
    key: "your-secret-token-here"
    name: "Admin Token"
    permissions:
      - upload
      - delete
      - list
```

## API Endpoints

### 1. Upload File

```bash
POST /upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-token" \
  -F "file=@photo.jpg" \
  -F "tag=images" \
  -F "public=true"
```

### 2. Download/View File

```bash
GET /:tag/:filename
Optional: ?token=xxx (for private files)
Optional: ?download=true (force download)

# Public file
curl http://localhost:8080/images/photo_abc123.jpg

# Private file with token
curl http://localhost:8080/documents/doc_xyz789.pdf?token=your-token
```

### 3. List Files

```bash
GET /api/files
Authorization: Bearer {token}
Optional: ?tag=images&page=1&limit=20&search=photo

curl -X GET "http://localhost:8080/api/files?tag=images&page=1&limit=20" \
  -H "Authorization: Bearer your-token"
```

### 4. Delete File

```bash
DELETE /api/files/:tag/:filename
Authorization: Bearer {token}

curl -X DELETE http://localhost:8080/api/files/images/photo_abc123.jpg \
  -H "Authorization: Bearer your-token"
```

### 5. Health Check

```bash
GET /health

curl http://localhost:8080/health
```

## Development

### Available Make Commands

```bash
make run          # Run the application
make build        # Build binary
make build-linux  # Build for Linux
make test         # Run tests
make clean        # Clean build artifacts
make fmt          # Format code
```

### Running Tests

```bash
make test
# or
go test -v ./...
```

## Deployment

### Build for Production

```bash
make build-linux
```

### Systemd Service

Create `/etc/systemd/system/cdn-fileserver.service`:

```ini
[Unit]
Description=CDN File Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/cdn-fileserver
ExecStart=/opt/cdn-fileserver/bin/cdn-fileserver
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name cdn.maarifnu.or.id;

    client_max_body_size 50M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Security

- Use strong random tokens (minimum 32 characters)
- Never commit `config.yaml` to version control
- Always use HTTPS in production
- Regularly update dependencies
- Enable file content validation
- Implement rate limiting (optional)

## License

MIT License

## Author

Development Team - MA'ARIF NU

---

**Base URL**:
- Development: http://localhost:8080
- Production: https://cdn.maarifnu.or.id
