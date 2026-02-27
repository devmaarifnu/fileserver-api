# TODO: CDN File Server API

## 📋 Overview
Dokumen ini berisi requirement lengkap untuk **CDN File Server API** yang digunakan untuk mengelola file upload, download, list, dan delete dengan sistem tag-based organization. API ini mendukung autentikasi token dan kontrol akses public/private untuk setiap file.

---

## 🎯 Tech Stack

### Backend Framework
- **Golang** - Programming language
- **Gin Framework** - HTTP web framework
- **File-based Storage** - No database required

### Libraries & Tools
```go
github.com/gin-gonic/gin              // Web framework
github.com/gin-contrib/cors           // CORS middleware
github.com/google/uuid                // UUID generator untuk unique filename
github.com/spf13/viper                // Configuration management (YAML)
github.com/sirupsen/logrus            // Logging library
gopkg.in/natefinch/lumberjack.v2     // Log rotation
```

---

## 🎯 Base URL
```
Production: https://cdn.maarifnu.or.id
Development: http://localhost:8080
```

---

## 📌 API Characteristics

### Features:
- ✅ **Token-based Authentication** - Multiple tokens dengan permissions
- ✅ **File Upload** - Dengan tag-based organization
- ✅ **Public/Private Files** - Control akses per file
- ✅ **File Download/View** - Serve file langsung
- ✅ **File List** - List dengan filter dan pagination
- ✅ **File Delete** - Hapus file dengan authorization
- ✅ **CORS Enabled** - Allow frontend access
- ✅ **Logging** - Comprehensive logging untuk semua operations
- ✅ **Response Compression** - Gzip compression

### File Organization:
- 📁 Tag-based directory structure: `storage/{tag}/{filename}`
- 📄 Metadata storage: `{filename}.meta.json`
- 🔐 Permission-based access control
- 🔒 Unique filename generation untuk prevent collision

---

## 🗂️ API Endpoints

### 1. UPLOAD FILE

#### 1.1 Upload Single File
```http
POST /upload
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Form Data:**
- `file` (required): File binary
- `tag` (required): Tag untuk organization (alphanumeric, dash, underscore)
- `public` (optional, default: false): Public atau private file

**Example Request:**
```bash
curl -X POST https://cdn.maarifnu.or.id/upload \
  -H "Authorization: Bearer your-secret-token" \
  -F "file=@photo.jpg" \
  -F "tag=images" \
  -F "public=true"
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "file_id": "photo_a1b2c3d4-e5f6-7890-abcd.jpg",
    "original_name": "photo.jpg",
    "url": "https://cdn.maarifnu.or.id/images/photo_a1b2c3d4-e5f6-7890-abcd.jpg",
    "tag": "images",
    "size": 1024576,
    "content_type": "image/jpeg",
    "public": true,
    "uploaded_at": "2025-01-16T10:30:00Z",
    "uploaded_by": "Admin Token"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Validation error",
  "errors": {
    "tag": "Tag is required and must be alphanumeric",
    "file": "File size exceeds maximum limit (50MB)"
  }
}
```

**Response Error (401):**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or missing token"
}
```

**Response Error (403):**
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Token does not have upload permission"
}
```

**Response Error (413):**
```json
{
  "success": false,
  "message": "File too large",
  "error": "Maximum file size is 50MB"
}
```

---

### 2. DOWNLOAD/VIEW FILE

#### 2.1 Get File by Tag and Filename
```http
GET /:tag/:filename
Optional Query: ?token=xxx (untuk private files)
Optional Query: ?download=true (force download)
```

**Path Parameters:**
- `tag` (required): File tag/category
- `filename` (required): Unique filename

**Example Requests:**
```bash
# Public file (no auth required)
curl https://cdn.maarifnu.or.id/images/photo_a1b2c3d4.jpg

# Private file dengan token di query
curl https://cdn.maarifnu.or.id/documents/report_x1y2z3.pdf?token=your-secret-token

# Private file dengan token di header
curl https://cdn.maarifnu.or.id/documents/report_x1y2z3.pdf \
  -H "Authorization: Bearer your-secret-token"

# Force download
curl https://cdn.maarifnu.or.id/images/photo_a1b2c3d4.jpg?download=true
```

**Response Success (200):**
- Returns file binary dengan proper Content-Type header
- Untuk download: Content-Disposition: attachment

**Response Error (403):**
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "This file is private and requires authentication"
}
```

**Response Error (404):**
```json
{
  "success": false,
  "message": "File not found",
  "error": "File with tag 'images' and filename 'photo.jpg' does not exist"
}
```

---

### 3. LIST FILES

#### 3.1 Get All Files with Filtering
```http
GET /api/files
Authorization: Bearer {token}
```

**Query Parameters:**
- `tag` (optional): Filter by specific tag
- `public` (optional): Filter by public status (true/false)
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 50, max: 100): Items per page
- `sort` (optional, default: desc): Sort order by upload date (asc/desc)
- `search` (optional): Search by filename

**Example Requests:**
```bash
# Get all files
curl -X GET https://cdn.maarifnu.or.id/api/files \
  -H "Authorization: Bearer your-secret-token"

# Filter by tag
curl -X GET "https://cdn.maarifnu.or.id/api/files?tag=images&page=1&limit=20" \
  -H "Authorization: Bearer your-secret-token"

# Filter by public status
curl -X GET "https://cdn.maarifnu.or.id/api/files?public=true" \
  -H "Authorization: Bearer your-secret-token"

# Search files
curl -X GET "https://cdn.maarifnu.or.id/api/files?search=photo" \
  -H "Authorization: Bearer your-secret-token"
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Files retrieved successfully",
  "data": {
    "files": [
      {
        "file_id": "photo_a1b2c3d4.jpg",
        "original_name": "photo.jpg",
        "tag": "images",
        "url": "https://cdn.maarifnu.or.id/images/photo_a1b2c3d4.jpg",
        "size": 1024576,
        "content_type": "image/jpeg",
        "public": true,
        "uploaded_at": "2025-01-16T10:30:00Z",
        "uploaded_by": "Admin Token"
      },
      {
        "file_id": "document_x1y2z3.pdf",
        "original_name": "report.pdf",
        "tag": "documents",
        "url": "https://cdn.maarifnu.or.id/documents/document_x1y2z3.pdf",
        "size": 2048000,
        "content_type": "application/pdf",
        "public": false,
        "uploaded_at": "2025-01-16T09:15:00Z",
        "uploaded_by": "Upload Token"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 95,
      "items_per_page": 20,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

**Response Error (401):**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or missing token"
}
```

**Response Error (403):**
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Token does not have list permission"
}
```

---

### 4. DELETE FILE

#### 4.1 Delete File by Tag and Filename
```http
DELETE /api/files/:tag/:filename
Authorization: Bearer {token}
```

**Path Parameters:**
- `tag` (required): File tag/category
- `filename` (required): Unique filename

**Example Request:**
```bash
curl -X DELETE https://cdn.maarifnu.or.id/api/files/images/photo_a1b2c3d4.jpg \
  -H "Authorization: Bearer your-secret-token"
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "File deleted successfully",
  "data": {
    "file_id": "photo_a1b2c3d4.jpg",
    "tag": "images",
    "deleted_at": "2025-01-16T11:45:00Z"
  }
}
```

**Response Error (401):**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or missing token"
}
```

**Response Error (403):**
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Token does not have delete permission"
}
```

**Response Error (404):**
```json
{
  "success": false,
  "message": "File not found",
  "error": "File with tag 'images' and filename 'photo.jpg' does not exist"
}
```

---

### 5. HEALTH CHECK

#### 5.1 Check API Health
```http
GET /health
```

**Response Success (200):**
```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "48h30m15s",
  "storage": {
    "total_files": 1250,
    "total_size": "2.5 GB",
    "disk_usage": "15%"
  }
}
```

---

## 🗂️ Project Structure

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
│   ├── images/
│   ├── documents/
│   ├── videos/
│   └── .gitkeep
├── logs/                          # Application logs
│   └── app.log
├── config.yaml                    # Configuration file
├── config.example.yaml            # Configuration example
├── .env                           # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                       # Build commands
└── README.md
```

---

## 🔧 Installation & Setup

### Prerequisites
```bash
- Go 1.21+
- Linux/Unix based OS (untuk production)
```

### Configuration File (config.yaml)
```yaml
# Application Configuration
app:
  name: "CDN File Server"
  env: "development"  # development, staging, production
  port: 8080
  version: "1.0.0"
  domain: "cdn.maarifnu.or.id"

# Storage Configuration
storage:
  base_path: "./storage"
  max_file_size: 52428800  # 50MB in bytes
  allowed_extensions:
    - jpg
    - jpeg
    - png
    - gif
    - webp
    - pdf
    - doc
    - docx
    - xls
    - xlsx
    - ppt
    - pptx
    - zip
    - rar
    - mp4
    - avi
    - mov

# Authentication Tokens
tokens:
  - id: "token_001"
    key: "your-secret-token-admin-here"
    name: "Admin Token"
    permissions:
      - upload
      - delete
      - list
  
  - id: "token_002"
    key: "your-secret-token-upload-here"
    name: "Upload Token"
    permissions:
      - upload
      - list
  
  - id: "token_003"
    key: "your-secret-token-readonly-here"
    name: "Readonly Token"
    permissions:
      - list

# CORS Configuration
cors:
  enabled: true
  allowed_origins:
    - "http://localhost:3000"
    - "https://maarifnu.or.id"
    - "https://www.maarifnu.or.id"
  allowed_methods:
    - "GET"
    - "POST"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Content-Type"
    - "Authorization"
  allow_credentials: true

# Logging Configuration
logging:
  level: "info"  # debug, info, warn, error
  format: "json" # json or text
  output: "file" # console, file, both
  file_path: "./logs/app.log"
  max_size: 100      # MB
  max_backups: 5
  max_age: 30        # days
  compress: true

# Security Configuration
security:
  validate_file_content: true  # Validate file magic bytes
  sanitize_filename: true      # Sanitize user input filename
```

### Installation Steps

1. **Clone Repository**
```bash
git clone https://github.com/maarifnu/cdn-fileserver.git
cd cdn-fileserver
```

2. **Install Dependencies**
```bash
go mod download
```

3. **Setup Configuration File**
```bash
cp config.example.yaml config.yaml
# Edit config.yaml dengan konfigurasi Anda
```

4. **Generate Strong Tokens**
```bash
# Generate random tokens untuk production
openssl rand -hex 32  # Generate 64 character hex token
```

5. **Create Storage Directory**
```bash
mkdir -p storage logs
chmod 755 storage logs
```

6. **Run Application**
```bash
# Development mode
go run cmd/api/main.go

# Build for production
go build -o bin/cdn-fileserver cmd/api/main.go

# Run production binary
./bin/cdn-fileserver
```

---

## 📝 Implementation TODO List

### Phase 1: Project Setup (Day 1)

#### 1.1 Initialize Project
- [ ] Create project directory: `mkdir cdn-fileserver && cd cdn-fileserver`
- [ ] Initialize Go module: `go mod init github.com/maarifnu/cdn-fileserver`
- [ ] Install dependencies:
  ```bash
  go get github.com/gin-gonic/gin
  go get github.com/google/uuid
  go get github.com/spf13/viper
  go get github.com/sirupsen/logrus
  go get gopkg.in/natefinch/lumberjack.v2
  go get github.com/gin-contrib/cors
  ```

#### 1.2 Create Directory Structure
- [ ] Create `cmd/api/main.go`
- [ ] Create `internal/` subdirectories:
  - `config/`
  - `models/`
  - `handlers/`
  - `services/`
  - `middleware/`
  - `utils/`
  - `routes/`
- [ ] Create `pkg/logger/`
- [ ] Create `storage/` directory dengan `.gitkeep`
- [ ] Create `logs/` directory dengan `.gitkeep`

#### 1.3 Setup Basic Files
- [ ] Create `config.yaml`
- [ ] Create `config.example.yaml`
- [ ] Create `.gitignore`
- [ ] Create `.env` untuk environment variables
- [ ] Create `Makefile` untuk build commands
- [ ] Create `README.md`

---

### Phase 2: Configuration & Logger (Day 1-2)

#### 2.1 Configuration Module (internal/config/config.go)
- [ ] Define `Config` struct dengan semua config sections
- [ ] Define `AppConfig` struct
- [ ] Define `StorageConfig` struct
- [ ] Define `TokenConfig` struct
- [ ] Define `CORSConfig` struct
- [ ] Define `LoggingConfig` struct
- [ ] Define `SecurityConfig` struct
- [ ] Implement `Load()` function menggunakan Viper
- [ ] Implement config validation
- [ ] Add environment variables override support
- [ ] Handle config file not found error

#### 2.2 Logger Module (pkg/logger/logger.go)
- [ ] Setup logrus configuration
- [ ] Implement log rotation dengan lumberjack
- [ ] Support JSON dan text format
- [ ] Support multiple output (console, file, both)
- [ ] Implement structured logging fields
- [ ] Add request ID untuk tracing
- [ ] Create helper functions: `Info()`, `Error()`, `Warn()`, `Debug()`

---

### Phase 3: Models & Utils (Day 2)

#### 3.1 File Metadata Model (internal/models/file_meta.go)
- [ ] Define `FileMeta` struct:
  ```go
  type FileMeta struct {
      FileID       string    `json:"file_id"`
      OriginalName string    `json:"original_name"`
      Tag          string    `json:"tag"`
      Size         int64     `json:"size"`
      ContentType  string    `json:"content_type"`
      Public       bool      `json:"public"`
      UploadedAt   time.Time `json:"uploaded_at"`
      UploadedBy   string    `json:"uploaded_by"`
  }
  ```
- [ ] Implement `Save()` method - save metadata to JSON file
- [ ] Implement `Load()` method - load metadata from JSON file
- [ ] Implement `Delete()` method - delete metadata file
- [ ] Implement `Exists()` method - check if metadata exists

#### 3.2 Response Utils (internal/utils/response.go)
- [ ] Implement `SuccessResponse()` function
- [ ] Implement `ErrorResponse()` function
- [ ] Implement `ValidationErrorResponse()` function
- [ ] Implement `PaginationResponse()` function
- [ ] Define response structs: `APIResponse`, `PaginationMeta`

#### 3.3 Validator Utils (internal/utils/validator.go)
- [ ] Implement `ValidateTag()` - alphanumeric, dash, underscore only
- [ ] Implement `ValidateFileExtension()` - check allowed extensions
- [ ] Implement `ValidateFileSize()` - check max file size
- [ ] Implement `SanitizeFilename()` - sanitize user input
- [ ] Implement `ValidateFileMagicBytes()` - validate file content (optional)

#### 3.4 Filename Utils (internal/utils/filename.go)
- [ ] Implement `GenerateUniqueFilename()` - generate filename dengan UUID
- [ ] Implement `ExtractExtension()` - extract file extension
- [ ] Implement `SanitizeName()` - clean filename
- [ ] Format: `{sanitized_name}_{uuid}.{ext}`

#### 3.5 File Utils (internal/utils/file.go)
- [ ] Implement `GetContentType()` - detect MIME type
- [ ] Implement `CreateDirectory()` - create directory if not exists
- [ ] Implement `GetFileSize()` - get file size
- [ ] Implement `FileExists()` - check file exists

---

### Phase 4: Middleware (Day 3)

#### 4.1 Auth Middleware (internal/middleware/auth.go)
- [ ] Implement `TokenAuth()` middleware function
- [ ] Extract token dari header `Authorization: Bearer {token}`
- [ ] Extract token dari query param `?token=xxx` (fallback)
- [ ] Validate token terhadap config tokens
- [ ] Check token permissions
- [ ] Set token info ke Gin context
- [ ] Return 401 untuk invalid token
- [ ] Return 403 untuk insufficient permissions
- [ ] Log authentication attempts

#### 4.2 Logger Middleware (internal/middleware/logger.go)
- [ ] Implement `LoggerMiddleware()` function
- [ ] Generate unique request ID
- [ ] Log incoming request: method, path, IP, user-agent
- [ ] Log response: status code, duration, response size
- [ ] Include request ID dalam semua logs
- [ ] Handle panic dan log error dengan stack trace

#### 4.3 CORS Middleware (internal/middleware/cors.go)
- [ ] Implement `CORSMiddleware()` function
- [ ] Load CORS config dari config.yaml
- [ ] Set proper CORS headers
- [ ] Handle preflight OPTIONS requests
- [ ] Use gin-contrib/cors library

#### 4.4 Recovery Middleware (internal/middleware/recovery.go)
- [ ] Implement `RecoveryMiddleware()` function
- [ ] Catch panic
- [ ] Log error dengan stack trace
- [ ] Return 500 Internal Server Error
- [ ] Prevent application crash

---

### Phase 5: Services (Day 3-4)

#### 5.1 Storage Service (internal/services/storage_service.go)
- [ ] Define `StorageService` interface
- [ ] Implement `SaveFile()` - save uploaded file
- [ ] Implement `GetFile()` - retrieve file
- [ ] Implement `DeleteFile()` - delete file
- [ ] Implement `ListFiles()` - list all files dengan filter
- [ ] Implement `GetStorageInfo()` - get storage statistics
- [ ] Handle directory creation
- [ ] Handle file permissions (0644)

#### 5.2 File Service (internal/services/file_service.go)
- [ ] Define `FileService` interface
- [ ] Implement `Upload()` - handle upload flow:
  - Validate input
  - Generate unique filename
  - Save file
  - Create metadata
  - Return file info
- [ ] Implement `Download()` - handle download flow:
  - Check file exists
  - Load metadata
  - Check access permission (public/private)
  - Return file
- [ ] Implement `List()` - handle list flow:
  - Scan storage directory
  - Load metadata for each file
  - Apply filters
  - Apply pagination
  - Return files list
- [ ] Implement `Delete()` - handle delete flow:
  - Check file exists
  - Delete physical file
  - Delete metadata
  - Return success

---

### Phase 6: Handlers (Day 4-5)

#### 6.1 Upload Handler (internal/handlers/upload.go)
- [ ] Implement `UploadFile()` handler
- [ ] Parse multipart form data
- [ ] Extract file, tag, public parameters
- [ ] Validate inputs
- [ ] Check file size
- [ ] Check file extension
- [ ] Call FileService.Upload()
- [ ] Return JSON response dengan file URL
- [ ] Handle errors properly
- [ ] Log upload activity

#### 6.2 Download Handler (internal/handlers/download.go)
- [ ] Implement `DownloadFile()` handler
- [ ] Parse tag dan filename dari URL path
- [ ] Load file metadata
- [ ] Check if file is public
- [ ] If private, validate token (header atau query param)
- [ ] Set Content-Type header
- [ ] Handle `?download=true` query - set Content-Disposition
- [ ] Serve file dengan `c.File()`
- [ ] Log access activity
- [ ] Handle 404 Not Found
- [ ] Handle 403 Forbidden

#### 6.3 List Handler (internal/handlers/list.go)
- [ ] Implement `ListFiles()` handler
- [ ] Parse query parameters: tag, public, page, limit, sort, search
- [ ] Validate pagination parameters
- [ ] Call FileService.List()
- [ ] Build pagination metadata
- [ ] Return JSON response dengan files array
- [ ] Handle errors
- [ ] Log list operation

#### 6.4 Delete Handler (internal/handlers/delete.go)
- [ ] Implement `DeleteFile()` handler
- [ ] Parse tag dan filename dari URL path
- [ ] Check file exists
- [ ] Call FileService.Delete()
- [ ] Return success response
- [ ] Log delete operation dengan detail
- [ ] Handle 404 Not Found

#### 6.5 Health Check Handler (internal/handlers/health.go)
- [ ] Implement `HealthCheck()` handler
- [ ] Get application uptime
- [ ] Get storage statistics:
  - Total files count
  - Total storage size
  - Disk usage percentage
- [ ] Return health status JSON
- [ ] Always return 200 OK

---

### Phase 7: Routes & Main Application (Day 5)

#### 7.1 Routes Configuration (internal/routes/routes.go)
- [ ] Create `SetupRoutes()` function
- [ ] Register middleware:
  - CORS middleware
  - Logger middleware
  - Recovery middleware
- [ ] Register public routes:
  - `GET /:tag/:filename` - Download/view file (conditional auth)
  - `GET /health` - Health check
- [ ] Register authenticated routes (API group):
  - `POST /upload` - Upload file (auth + upload permission)
  - `GET /api/files` - List files (auth + list permission)
  - `DELETE /api/files/:tag/:filename` - Delete file (auth + delete permission)

#### 7.2 Main Application (cmd/api/main.go)
- [ ] Load configuration dari config.yaml
- [ ] Initialize logger
- [ ] Setup Gin router
- [ ] Set Gin mode (debug/release)
- [ ] Setup routes
- [ ] Add security headers middleware
- [ ] Setup graceful shutdown
- [ ] Start HTTP server
- [ ] Handle shutdown signals (SIGINT, SIGTERM)
- [ ] Log server startup
- [ ] Log server shutdown

---

### Phase 8: Security Enhancements (Day 6)

#### 8.1 Security Headers
- [ ] Add `X-Content-Type-Options: nosniff`
- [ ] Add `X-Frame-Options: DENY`
- [ ] Add `X-XSS-Protection: 1; mode=block`
- [ ] Add `Strict-Transport-Security` untuk HTTPS
- [ ] Remove `X-Powered-By` header

#### 8.2 Input Validation & Sanitization
- [ ] Sanitize all user inputs (tag, filename)
- [ ] Validate file magic bytes (optional)
- [ ] Prevent directory traversal attacks
- [ ] Validate Content-Type
- [ ] Check for malicious files

#### 8.3 Rate Limiting (Optional)
- [ ] Implement rate limiter middleware
- [ ] Limit: 100 requests per minute per IP
- [ ] Return 429 Too Many Requests
- [ ] Log rate limit violations

---

### Phase 9: Testing (Day 6-7)

#### 9.1 Manual Testing Checklist
- [ ] Upload file dengan tag "images"
- [ ] Upload file dengan tag "documents"
- [ ] Upload public file (`public=true`)
- [ ] Upload private file (`public=false`)
- [ ] Upload file dengan nama sama di tag berbeda ✓
- [ ] Upload file dengan nama sama di tag sama (harus unique) ✓
- [ ] Download public file tanpa token ✓
- [ ] Download private file tanpa token (harus 403) ✓
- [ ] Download private file dengan token valid ✓
- [ ] Download private file dengan token di query param
- [ ] Download file dengan `?download=true`
- [ ] List all files dengan pagination
- [ ] List files filter by tag
- [ ] List files filter by public status
- [ ] Search files by filename
- [ ] Delete file dengan token valid
- [ ] Delete file dengan token invalid (harus 403)
- [ ] Delete file dengan insufficient permission (harus 403)
- [ ] Upload file melebihi max size (harus 413)
- [ ] Upload file dengan extension tidak diizinkan
- [ ] Test dengan invalid tag format
- [ ] Test dengan invalid token
- [ ] Test health check endpoint

#### 9.2 Test dengan cURL
```bash
# Upload public file
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer admin-token" \
  -F "file=@photo.jpg" \
  -F "tag=images" \
  -F "public=true"

# Upload private file
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer admin-token" \
  -F "file=@document.pdf" \
  -F "tag=documents" \
  -F "public=false"

# Download public file
curl http://localhost:8080/images/photo_uuid.jpg -o downloaded.jpg

# Download private file dengan token
curl http://localhost:8080/documents/doc_uuid.pdf?token=admin-token -o doc.pdf

# List all files
curl -X GET http://localhost:8080/api/files \
  -H "Authorization: Bearer admin-token"

# List files by tag
curl -X GET "http://localhost:8080/api/files?tag=images&page=1&limit=20" \
  -H "Authorization: Bearer admin-token"

# Search files
curl -X GET "http://localhost:8080/api/files?search=photo" \
  -H "Authorization: Bearer admin-token"

# Delete file
curl -X DELETE http://localhost:8080/api/files/images/photo_uuid.jpg \
  -H "Authorization: Bearer admin-token"

# Health check
curl http://localhost:8080/health
```

#### 9.3 Unit Tests (Optional)
- [ ] Test filename generator
- [ ] Test validators
- [ ] Test authentication middleware
- [ ] Test response helpers
- [ ] Test config loader

---

### Phase 10: Documentation (Day 7)

#### 10.1 README.md
- [ ] Project description
- [ ] Features list
- [ ] Tech stack
- [ ] Prerequisites
- [ ] Installation steps
- [ ] Configuration guide
- [ ] Usage examples
- [ ] API documentation
- [ ] Deployment guide
- [ ] Troubleshooting

#### 10.2 API Documentation
- [ ] Endpoint specifications
- [ ] Request/response examples
- [ ] Authentication guide
- [ ] Error codes reference
- [ ] Query parameters documentation
- [ ] cURL examples

#### 10.3 Configuration Documentation
- [ ] config.yaml explanation
- [ ] Token generation guide
- [ ] Security best practices
- [ ] Performance tuning tips

---

### Phase 11: Deployment Preparation (Day 8)

#### 11.1 Build for Production
- [ ] Create production `config.yaml`
- [ ] Generate strong random tokens (32+ characters)
- [ ] Build binary: `go build -o bin/cdn-fileserver cmd/api/main.go`
- [ ] Test production binary
- [ ] Create deployment checklist

#### 11.2 Systemd Service
- [ ] Create systemd service file:
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
  StandardOutput=journal
  StandardError=journal

  [Install]
  WantedBy=multi-user.target
  ```
- [ ] Enable service: `sudo systemctl enable cdn-fileserver`
- [ ] Start service: `sudo systemctl start cdn-fileserver`
- [ ] Check status: `sudo systemctl status cdn-fileserver`

#### 11.3 Nginx Reverse Proxy
- [ ] Create Nginx config file `/etc/nginx/sites-available/cdn.maarifnu.or.id`:
  ```nginx
  # HTTP redirect to HTTPS
  server {
      listen 80;
      server_name cdn.maarifnu.or.id;
      return 301 https://$server_name$request_uri;
  }

  # HTTPS server
  server {
      listen 443 ssl http2;
      server_name cdn.maarifnu.or.id;

      # SSL Configuration
      ssl_certificate /etc/letsencrypt/live/cdn.maarifnu.or.id/fullchain.pem;
      ssl_certificate_key /etc/letsencrypt/live/cdn.maarifnu.or.id/privkey.pem;
      ssl_protocols TLSv1.2 TLSv1.3;
      ssl_ciphers HIGH:!aNULL:!MD5;
      ssl_prefer_server_ciphers on;

      # Upload size limit
      client_max_body_size 50M;
      client_body_buffer_size 50M;

      # Security headers
      add_header X-Frame-Options "DENY" always;
      add_header X-Content-Type-Options "nosniff" always;
      add_header X-XSS-Protection "1; mode=block" always;
      add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

      # Gzip compression
      gzip on;
      gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

      # Cache static files
      location ~* \.(jpg|jpeg|png|gif|ico|css|js|webp|svg)$ {
          expires 1y;
          add_header Cache-Control "public, immutable";
          proxy_pass http://localhost:8080;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto $scheme;
      }

      # API endpoints
      location / {
          proxy_pass http://localhost:8080;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto $scheme;
          
          # Timeout untuk upload besar
          proxy_connect_timeout 600;
          proxy_send_timeout 600;
          proxy_read_timeout 600;
          send_timeout 600;
      }
  }
  ```
- [ ] Enable site: `sudo ln -s /etc/nginx/sites-available/cdn.maarifnu.or.id /etc/nginx/sites-enabled/`
- [ ] Test config: `sudo nginx -t`
- [ ] Reload Nginx: `sudo systemctl reload nginx`

#### 11.4 SSL Certificate (Let's Encrypt)
- [ ] Install Certbot: `sudo apt install certbot python3-certbot-nginx`
- [ ] Generate certificate: `sudo certbot --nginx -d cdn.maarifnu.or.id`
- [ ] Test auto-renewal: `sudo certbot renew --dry-run`
- [ ] Setup cron for auto-renewal

#### 11.5 File Permissions
- [ ] Set ownership: `sudo chown -R www-data:www-data /opt/cdn-fileserver`
- [ ] Set directory permissions: `sudo chmod -R 755 /opt/cdn-fileserver`
- [ ] Set storage permissions: `sudo chmod -R 755 /opt/cdn-fileserver/storage`
- [ ] Set logs permissions: `sudo chmod -R 755 /opt/cdn-fileserver/logs`
- [ ] Set file permissions: `sudo chmod 644 /opt/cdn-fileserver/config.yaml`
- [ ] Protect sensitive files: `sudo chmod 600 /opt/cdn-fileserver/.env`

---

### Phase 12: Monitoring & Maintenance (Day 8-9)

#### 12.1 Log Monitoring
- [ ] Setup logrotate untuk application logs:
  ```
  /opt/cdn-fileserver/logs/*.log {
      daily
      rotate 30
      compress
      delaycompress
      notifempty
      create 0644 www-data www-data
      sharedscripts
      postrotate
          systemctl reload cdn-fileserver > /dev/null 2>&1 || true
      endscript
  }
  ```
- [ ] Monitor logs: `sudo tail -f /opt/cdn-fileserver/logs/app.log`
- [ ] Setup alerts untuk critical errors
- [ ] Use log aggregation tool (optional: ELK stack, Grafana Loki)

#### 12.2 Disk Space Monitoring
- [ ] Monitor storage usage: `du -sh /opt/cdn-fileserver/storage/*`
- [ ] Setup alerts jika disk > 80%
- [ ] Implement cleanup script untuk old files (optional)
- [ ] Plan storage expansion strategy

#### 12.3 Application Monitoring
- [ ] Check service status: `sudo systemctl status cdn-fileserver`
- [ ] Monitor memory usage: `ps aux | grep cdn-fileserver`
- [ ] Monitor CPU usage
- [ ] Check open connections: `sudo netstat -tulpn | grep 8080`
- [ ] Monitor response time
- [ ] Setup uptime monitoring (optional: UptimeRobot, Pingdom)

#### 12.4 Backup Strategy
- [ ] Backup storage directory daily
- [ ] Backup config files
- [ ] Test restore procedures
- [ ] Keep offsite backups
- [ ] Document recovery procedures
- [ ] Automate backup with cron:
  ```bash
  # Backup storage daily at 2 AM
  0 2 * * * tar -czf /backup/cdn-storage-$(date +\%Y\%m\%d).tar.gz /opt/cdn-fileserver/storage
  ```

---

## 🎉 Launch Checklist

### Pre-Launch
- [ ] All endpoints tested ✓
- [ ] Security measures implemented ✓
- [ ] Logging berjalan normal ✓
- [ ] SSL certificate installed ✓
- [ ] Nginx reverse proxy configured ✓
- [ ] Systemd service running ✓
- [ ] File permissions correct ✓
- [ ] Documentation complete ✓
- [ ] Backup strategy in place ✓
- [ ] Monitoring configured ✓

### Production Testing
- [ ] Upload test file dari production
- [ ] Download public file
- [ ] Download private file dengan token
- [ ] Test file list API
- [ ] Test file delete API
- [ ] Verify logs are working
- [ ] Check disk space
- [ ] Verify HTTPS working
- [ ] Test CORS from frontend
- [ ] Performance test dengan load testing tool
- [ ] Security scan (optional: OWASP ZAP)

### Post-Launch
- [ ] Monitor error logs untuk 24 jam pertama
- [ ] Check application uptime
- [ ] Verify backup running
- [ ] Monitor disk usage trend
- [ ] Collect performance metrics
- [ ] Update documentation jika ada changes
- [ ] Share API documentation dengan team

---

## 🔮 Future Enhancements (Optional)

### Image Processing
- [ ] Image resize on-the-fly: `?w=800&h=600`
- [ ] Image quality control: `?q=80`
- [ ] Automatic thumbnail generation
- [ ] WebP conversion untuk browser support
- [ ] Format conversion (jpg to png, etc)

### Video Processing
- [ ] Video transcoding
- [ ] Generate video thumbnails
- [ ] Stream video dengan HLS/DASH

### Advanced Features
- [ ] File versioning
- [ ] Batch upload API
- [ ] Batch delete API
- [ ] File search dengan Elasticsearch
- [ ] Storage quota per token
- [ ] Usage analytics dashboard
- [ ] CDN integration (CloudFlare, AWS CloudFront)
- [ ] S3-compatible API
- [ ] Webhook notifications untuk file events
- [ ] Admin web interface untuk file management
- [ ] File compression untuk non-image files
- [ ] Automatic cleanup untuk expired files
- [ ] Multi-region replication

### Integration
- [ ] Integrate dengan external storage (S3, MinIO)
- [ ] Integrate dengan image optimization service (TinyPNG)
- [ ] Integrate dengan virus scanning (ClamAV)
- [ ] Integrate dengan monitoring (Prometheus, Grafana)

---

## 📅 Development Timeline

### Week 1 (Day 1-5)
- **Day 1:** Project setup, configuration, logger
- **Day 2:** Models, utils, validators
- **Day 3:** Middleware, services (partial)
- **Day 4:** Services (complete), handlers (partial)
- **Day 5:** Handlers (complete), routes, main application

### Week 2 (Day 6-9)
- **Day 6:** Security enhancements, testing
- **Day 7:** More testing, documentation
- **Day 8:** Deployment preparation, production setup
- **Day 9:** Monitoring, final testing, launch

### Total Estimated Time: 9 Working Days

---

## 📊 Performance Optimization Tips

### 1. Efficient File Serving
- [ ] Use `c.File()` untuk serve file efficiently
- [ ] Enable sendfile di Nginx untuk zero-copy file serving
- [ ] Implement proper cache headers

### 2. Concurrent Operations
- [ ] Use goroutines untuk background tasks
- [ ] Implement worker pool untuk batch operations
- [ ] Use context untuk timeout control

### 3. Memory Management
- [ ] Stream large files instead of loading to memory
- [ ] Limit concurrent uploads
- [ ] Implement file size limits

### 4. Caching Strategy (Optional)
- [ ] Cache metadata untuk frequently accessed files
- [ ] Use Redis untuk distributed caching
- [ ] Implement cache warming for popular files

---

## 🔐 Security Best Practices

### Token Management
- [ ] Generate strong random tokens (minimum 32 characters)
- [ ] Never commit tokens to version control
- [ ] Use environment variables untuk sensitive data
- [ ] Rotate tokens regularly (every 90 days)
- [ ] Implement token expiration (optional)
- [ ] Log all token usage

### File Security
- [ ] Validate file content, not just extension
- [ ] Check for malicious files
- [ ] Prevent directory traversal
- [ ] Set proper file permissions (0644)
- [ ] Sanitize all user inputs
- [ ] Implement virus scanning (optional)

### Network Security
- [ ] Always use HTTPS in production
- [ ] Implement rate limiting
- [ ] Use WAF (Web Application Firewall) - optional
- [ ] Monitor for suspicious activities
- [ ] Keep SSL certificates updated
- [ ] Use strong SSL configuration

### Application Security
- [ ] Keep dependencies updated
- [ ] Regular security audits
- [ ] Implement proper error handling (don't leak info)
- [ ] Use secure headers
- [ ] Disable directory listing
- [ ] Remove debug information in production

---

## 📞 Support & Troubleshooting

### Common Issues

**Issue: "Failed to load configuration"**
- Check if config.yaml exists
- Verify YAML syntax
- Check file permissions

**Issue: "Permission denied" saat upload**
- Check storage directory permissions: `chmod 755 storage/`
- Check ownership: `chown www-data:www-data storage/`

**Issue: "File too large"**
- Check `max_file_size` in config.yaml
- Check Nginx `client_max_body_size`
- Check system disk space

**Issue: "Token authentication failed"**
- Verify token format in config.yaml
- Check Authorization header format: `Bearer {token}`
- Check token has required permissions

**Issue: "CORS error dari frontend"**
- Check CORS configuration in config.yaml
- Verify allowed_origins includes frontend URL
- Check preflight OPTIONS handling

### Debugging Commands
```bash
# Check service status
sudo systemctl status cdn-fileserver

# View logs
sudo journalctl -u cdn-fileserver -f

# View application logs
tail -f /opt/cdn-fileserver/logs/app.log

# Check disk space
df -h

# Check file permissions
ls -la /opt/cdn-fileserver/storage/

# Test Nginx config
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx

# Restart service
sudo systemctl restart cdn-fileserver
```

---

## 📝 Notes

### Development Tips
- Use `gin.SetMode(gin.ReleaseMode)` untuk production
- Always validate dan sanitize user input
- Use context timeout untuk file operations
- Implement proper error handling
- Keep sensitive data di config file (jangan commit)
- Use environment variables untuk production secrets
- Write clean, maintainable code
- Add comments untuk complex logic
- Follow Go best practices dan conventions

### Maintenance Schedule
- **Daily:** Monitor logs dan disk space
- **Weekly:** Review error logs, check backups
- **Monthly:** Update dependencies, rotate tokens
- **Quarterly:** Security audit, performance review
- **Yearly:** Major version update, architecture review

---

**Created:** January 16, 2025  
**Project:** CDN File Server API  
**Stack:** Go + Gin Framework  
**Version:** 1.0.0  
**Author:** Development Team

---

**Good luck with the development! 🚀**
