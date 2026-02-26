# CDN File Server - API Contract Documentation

## Base Information

- **Base URL (Development):** `http://localhost:8080`
- **Base URL (Production):** `https://cdn.maarifnu.or.id`
- **API Version:** 1.0.0
- **Content Type:** `application/json` (responses), `multipart/form-data` (upload)

---

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer {your-token}
```

For file downloads, token can also be passed as query parameter:
```
?token={your-token}
```

### Token Permissions

| Permission | Description |
|------------|-------------|
| `upload`   | Can upload files |
| `list`     | Can list files |
| `delete`   | Can delete files |

---

## Standard Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful message",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description"
}
```

### Validation Error Response
```json
{
  "success": false,
  "message": "Validation error",
  "errors": {
    "field_name": "Error description"
  }
}
```

---

## Endpoints

### 1. Health Check

Check API health and server status.

**Endpoint:** `GET /health`

**Authentication:** Not required

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:** `200 OK`
```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "48h30m15s",
  "storage": {
    "total_files": 1250,
    "total_size": "2.5 GB"
  }
}
```

---

### 2. Upload File

Upload a file with tag-based organization.

**Endpoint:** `POST /upload`

**Authentication:** Required (permission: `upload`)

**Content-Type:** `multipart/form-data`

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | File | Yes | Binary file to upload |
| `tag` | String | Yes | Tag for organization (alphanumeric, dash, underscore only) |
| `public` | Boolean | No | Public (true) or private (false) file. Default: false |

**Request Example:**
```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-token" \
  -F "file=@photo.jpg" \
  -F "tag=images" \
  -F "public=true"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "file_id": "photo_a1b2c3d4.jpg",
    "original_name": "photo.jpg",
    "url": "http://localhost:8080/images/photo_a1b2c3d4.jpg",
    "tag": "images",
    "size": 1024576,
    "content_type": "image/jpeg",
    "public": true,
    "uploaded_at": "2025-01-28T10:30:00Z",
    "uploaded_by": "Admin Token"
  }
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| `400 Bad Request` | Validation error (invalid tag, file extension not allowed, etc.) |
| `401 Unauthorized` | Invalid or missing token |
| `403 Forbidden` | Token doesn't have upload permission |
| `413 Payload Too Large` | File exceeds maximum size (50MB) |

**Validation Rules:**
- Tag: alphanumeric, dash, underscore only (max 50 chars)
- File size: Maximum 50MB (configurable)
- File extensions: Only allowed extensions from config

---

### 3. Download/View File

Download or view a file.

**Endpoint:** `GET /:tag/:filename`

**Authentication:** Optional (required for private files)

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `token` | String | No | Authentication token (alternative to header) |
| `download` | Boolean | No | Force download instead of inline view |

**Request Examples:**

1. **Public file (no auth):**
```bash
curl http://localhost:8080/images/photo_a1b2c3d4.jpg
```

2. **Private file with header auth:**
```bash
curl http://localhost:8080/documents/report_xyz.pdf \
  -H "Authorization: Bearer your-token"
```

3. **Private file with query auth:**
```bash
curl http://localhost:8080/documents/report_xyz.pdf?token=your-token
```

4. **Force download:**
```bash
curl http://localhost:8080/images/photo_a1b2c3d4.jpg?download=true
```

**Response:** `200 OK`
- Returns file binary with appropriate Content-Type header
- For download=true: includes `Content-Disposition: attachment` header

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| `403 Forbidden` | File is private and requires authentication |
| `404 Not Found` | File not found |

---

### 4. List Files

Retrieve a list of files with filtering and pagination.

**Endpoint:** `GET /api/files`

**Authentication:** Required (permission: `list`)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `tag` | String | No | - | Filter by specific tag |
| `public` | Boolean | No | - | Filter by public status (true/false) |
| `page` | Integer | No | 1 | Page number |
| `limit` | Integer | No | 50 | Items per page (max: 100) |
| `sort` | String | No | desc | Sort order by upload date (asc/desc) |
| `search` | String | No | - | Search by filename |

**Request Examples:**

1. **Get all files:**
```bash
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer your-token"
```

2. **Filter by tag with pagination:**
```bash
curl "http://localhost:8080/api/files?tag=images&page=1&limit=20" \
  -H "Authorization: Bearer your-token"
```

3. **Search files:**
```bash
curl "http://localhost:8080/api/files?search=photo" \
  -H "Authorization: Bearer your-token"
```

4. **Filter public files only:**
```bash
curl "http://localhost:8080/api/files?public=true" \
  -H "Authorization: Bearer your-token"
```

**Response:** `200 OK`
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
        "url": "http://localhost:8080/images/photo_a1b2c3d4.jpg",
        "size": 1024576,
        "content_type": "image/jpeg",
        "public": true,
        "uploaded_at": "2025-01-28T10:30:00Z",
        "uploaded_by": "Admin Token"
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

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| `401 Unauthorized` | Invalid or missing token |
| `403 Forbidden` | Token doesn't have list permission |

---

### 5. Delete File

Delete a file and its metadata.

**Endpoint:** `DELETE /api/files/:tag/:filename`

**Authentication:** Required (permission: `delete`)

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tag` | String | Yes | File tag/category |
| `filename` | String | Yes | Unique filename |

**Request Example:**
```bash
curl -X DELETE http://localhost:8080/api/files/images/photo_a1b2c3d4.jpg \
  -H "Authorization: Bearer your-token"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "File deleted successfully",
  "data": {
    "file_id": "photo_a1b2c3d4.jpg",
    "tag": "images",
    "deleted_at": "2025-01-28T11:45:00Z"
  }
}
```

**Error Responses:**

| Status Code | Description |
|-------------|-------------|
| `401 Unauthorized` | Invalid or missing token |
| `403 Forbidden` | Token doesn't have delete permission |
| `404 Not Found` | File not found |

---

## HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| `200 OK` | Request successful |
| `400 Bad Request` | Validation error or malformed request |
| `401 Unauthorized` | Authentication required or invalid token |
| `403 Forbidden` | Access denied (insufficient permissions or private file) |
| `404 Not Found` | Resource not found |
| `413 Payload Too Large` | File size exceeds maximum limit |
| `500 Internal Server Error` | Server error |

---

## File Organization

### Storage Structure
```
storage/
├── {tag}/
│   ├── {filename}              # Actual file
│   └── {filename}.meta.json    # File metadata
```

### Filename Format
Generated filenames use the format:
```
{sanitized_original_name}_{uuid}.{extension}
```

Example: `photo_a1b2c3d4.jpg`

### Metadata Format
Each file has an associated `.meta.json` file:
```json
{
  "file_id": "photo_a1b2c3d4.jpg",
  "original_name": "photo.jpg",
  "tag": "images",
  "size": 1024576,
  "content_type": "image/jpeg",
  "public": true,
  "uploaded_at": "2025-01-28T10:30:00Z",
  "uploaded_by": "Admin Token"
}
```

---

## Rate Limiting

Currently not implemented, but recommended for production:
- 100 requests per minute per IP
- Returns `429 Too Many Requests` when exceeded

---

## Security Headers

All responses include security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`

Public files also include cache headers:
- `Cache-Control: public, max-age=31536000, immutable`

---

## Error Examples

### 1. Validation Error
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

### 2. Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Invalid or missing token"
}
```

### 3. Forbidden
```json
{
  "success": false,
  "message": "Forbidden",
  "error": "Token does not have delete permission"
}
```

### 4. Not Found
```json
{
  "success": false,
  "message": "File not found",
  "error": "File with tag 'images' and filename 'photo.jpg' does not exist"
}
```

---

## Configuration Limits

| Setting | Default Value | Configurable |
|---------|---------------|--------------|
| Max File Size | 50MB | Yes (config.yaml) |
| Max Items Per Page | 100 | Yes (hardcoded) |
| Default Items Per Page | 50 | Yes (configurable) |
| Request Timeout | 10 minutes | Yes (main.go) |

---

## Allowed File Extensions (Default)

### Images
- jpg, jpeg, png, gif, webp

### Documents
- pdf, doc, docx, xls, xlsx, ppt, pptx

### Archives
- zip, rar

### Videos
- mp4, avi, mov

*Note: Extensions are configurable in config.yaml*

---

## Best Practices

1. **Authentication**
   - Always use HTTPS in production
   - Keep tokens secure and never commit to version control
   - Rotate tokens regularly (every 90 days)

2. **File Upload**
   - Validate file types on client-side before upload
   - Use descriptive tags for better organization
   - Set appropriate public/private flags

3. **File Access**
   - Use query parameter token for embeddable resources
   - Use header token for API integrations
   - Cache public file URLs when possible

4. **Error Handling**
   - Always check response status code
   - Parse error messages for user feedback
   - Implement retry logic for 5xx errors

---

## CORS Configuration

CORS is enabled and configurable. Default allowed origins:
- `http://localhost:3000`
- `https://maarifnu.or.id`
- `https://www.maarifnu.or.id`

Allowed methods:
- GET, POST, DELETE, OPTIONS

---

## Webhook Support

Not currently implemented. Future enhancement planned.

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-01-28 | Initial release |

---

## Support

For issues or questions, please contact:
- Email: support@maarifnu.or.id
- GitHub: https://github.com/maarifnu/cdn-fileserver

---

**Last Updated:** January 28, 2025
