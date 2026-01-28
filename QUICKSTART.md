# Quick Start Guide - CDN File Server API

## 🚀 Get Started in 5 Minutes

### Prerequisites
- Go 1.21 or higher installed
- Git installed

### Step 1: Clone and Setup (1 minute)

```bash
cd cdn-fileserver
go mod download
cp config.example.yaml config.yaml
```

### Step 2: Configure (1 minute)

Edit `config.yaml` - at minimum, change the tokens:

```yaml
tokens:
  - id: "token_001"
    key: "YOUR-SECURE-TOKEN-HERE"  # Change this!
    name: "Admin Token"
    permissions:
      - upload
      - delete
      - list
```

> 💡 **Tip:** Generate secure token with: `openssl rand -hex 32`

### Step 3: Run (1 minute)

```bash
# Option 1: Using Make
make run

# Option 2: Using Go directly
go run cmd/api/main.go

# Option 3: Using built binary
make build
./bin/cdn-fileserver
```

You should see:
```
Server listening on :8080
Base URL: http://localhost:8080
```

### Step 4: Test (2 minutes)

#### 4.1 Health Check
```bash
curl http://localhost:8080/health
```

Expected output:
```json
{
  "status": "ok",
  "version": "1.0.0",
  ...
}
```

#### 4.2 Upload a File
```bash
# Create a test file
echo "Hello CDN" > test.pdf

# Upload it
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR-SECURE-TOKEN-HERE" \
  -F "file=@test.pdf" \
  -F "tag=documents" \
  -F "public=true"
```

You'll get a response with the file URL:
```json
{
  "success": true,
  "data": {
    "url": "http://localhost:8080/documents/test_abc123.pdf",
    ...
  }
}
```

#### 4.3 Access Your File
```bash
# Open in browser or use curl
curl http://localhost:8080/documents/test_abc123.pdf
```

---

## 🎯 Common Use Cases

### Use Case 1: Upload Image for Website

```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR-TOKEN" \
  -F "file=@logo.png" \
  -F "tag=images" \
  -F "public=true"
```

Then use the returned URL in your HTML:
```html
<img src="http://localhost:8080/images/logo_abc123.png" />
```

### Use Case 2: Private Document Storage

```bash
# Upload private document
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR-TOKEN" \
  -F "file=@confidential.pdf" \
  -F "tag=documents" \
  -F "public=false"

# Access with token
curl "http://localhost:8080/documents/confidential_xyz.pdf?token=YOUR-TOKEN"
```

### Use Case 3: List and Manage Files

```bash
# List all files
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer YOUR-TOKEN"

# Search for specific files
curl "http://localhost:8080/api/files?search=logo" \
  -H "Authorization: Bearer YOUR-TOKEN"

# Delete a file
curl -X DELETE http://localhost:8080/api/files/images/old_logo.png \
  -H "Authorization: Bearer YOUR-TOKEN"
```

---

## 🔑 Token Permissions Explained

### Admin Token (Full Access)
```yaml
permissions:
  - upload    # Can upload files
  - delete    # Can delete files
  - list      # Can list files
```
**Use for:** Admin panel, full management

### Upload Token (Upload + View)
```yaml
permissions:
  - upload    # Can upload files
  - list      # Can list files
```
**Use for:** User file uploads, content management

### Readonly Token (View Only)
```yaml
permissions:
  - list      # Can only list files
```
**Use for:** Public viewing, reporting, analytics

---

## 📁 File Organization Tips

### Use Descriptive Tags
```bash
# Good - organized by type
tag=profile-pictures
tag=product-images
tag=invoices-2025
tag=user-documents

# Bad - too generic
tag=files
tag=uploads
```

### Tag Naming Rules
- ✅ Alphanumeric characters
- ✅ Dashes and underscores
- ✅ Max 50 characters
- ❌ No spaces
- ❌ No special characters

---

## 🔒 Security Best Practices

### 1. Generate Strong Tokens
```bash
# Use this command to generate secure tokens
openssl rand -hex 32

# Output example:
# a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456
```

### 2. Never Commit Tokens
```bash
# Already in .gitignore
echo "config.yaml" >> .gitignore
git add .gitignore
```

### 3. Use HTTPS in Production
```yaml
# In production config
app:
  env: "production"
  domain: "cdn.yourdomain.com"  # Will use https://
```

### 4. Set Appropriate File Permissions
```bash
# On Linux/Unix
chmod 600 config.yaml          # Only owner can read/write
chmod 755 storage/             # Owner full, others read/execute
```

---

## 🛠️ Troubleshooting

### Problem: "Failed to load configuration"
**Solution:** Make sure `config.yaml` exists in the project root
```bash
ls -la config.yaml
# If not found:
cp config.example.yaml config.yaml
```

### Problem: "Permission denied" when uploading
**Solution:** Check storage directory permissions
```bash
chmod 755 storage/
chmod 755 logs/
```

### Problem: "Invalid or missing token"
**Solution:** Check your token in the request
```bash
# Make sure Authorization header is correct
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer YOUR-ACTUAL-TOKEN"
#                             ^^^ Space is important!
```

### Problem: "File extension not allowed"
**Solution:** Add the extension to config.yaml
```yaml
storage:
  allowed_extensions:
    - jpg
    - png
    - pdf
    - txt  # Add your extension here
```

### Problem: Port 8080 already in use
**Solution:** Change the port in config.yaml
```yaml
app:
  port: 8081  # Change to any available port
```

---

## 📊 Monitoring

### Check Server Status
```bash
curl http://localhost:8080/health
```

### View Logs
```bash
# Console output (if running in terminal)
# Check logs/ directory

# Tail log file
tail -f logs/app.log

# On Windows with Git Bash
tail -f logs/app.log

# View last 100 lines
tail -n 100 logs/app.log
```

### Storage Usage
The health check endpoint shows storage statistics:
```json
{
  "storage": {
    "total_files": 1250,
    "total_size": "2.5 GB"
  }
}
```

---

## 🚀 Production Deployment

### Quick Production Checklist

1. **Generate Strong Tokens**
   ```bash
   openssl rand -hex 32  # Run 3 times for 3 tokens
   ```

2. **Update Config for Production**
   ```yaml
   app:
     env: "production"
     domain: "cdn.yourdomain.com"

   logging:
     level: "info"  # Not debug
     format: "json"
   ```

3. **Build Binary**
   ```bash
   make build-linux  # For Linux server
   # or
   make build        # For current OS
   ```

4. **Setup Systemd Service** (Linux)
   ```bash
   sudo nano /etc/systemd/system/cdn-fileserver.service
   # Add service configuration (see README.md)

   sudo systemctl enable cdn-fileserver
   sudo systemctl start cdn-fileserver
   ```

5. **Setup Nginx Reverse Proxy**
   ```nginx
   server {
       listen 80;
       server_name cdn.yourdomain.com;

       client_max_body_size 50M;

       location / {
           proxy_pass http://localhost:8080;
       }
   }
   ```

6. **Setup SSL with Let's Encrypt**
   ```bash
   sudo certbot --nginx -d cdn.yourdomain.com
   ```

See [README.md](README.md) for detailed deployment guide.

---

## 📚 Next Steps

- Read the full [API Contract Documentation](API-CONTRACT.md)
- Import [Postman Collection](CDN-FileServer.postman_collection.json)
- Check the detailed [README](README.md)
- Review the [TODO](TODO-CDN-FILE-SERVER.md) for advanced features

---

## 💡 Tips & Tricks

### Batch Upload with Script
```bash
#!/bin/bash
TOKEN="your-token-here"

for file in images/*.jpg; do
  echo "Uploading $file..."
  curl -X POST http://localhost:8080/upload \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@$file" \
    -F "tag=gallery" \
    -F "public=true"
done
```

### Download with Original Filename
```bash
# Use ?download=true to force download
curl "http://localhost:8080/images/photo_abc.jpg?download=true" -o photo.jpg
```

### List Files as JSON
```bash
# Pretty print with jq
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer TOKEN" | jq .

# Save to file
curl http://localhost:8080/api/files \
  -H "Authorization: Bearer TOKEN" > files.json
```

---

## 🎓 Learning Resources

### Project Structure
```
cdn-fileserver/
├── cmd/api/main.go           # Start here - application entry
├── internal/
│   ├── config/               # Configuration management
│   ├── handlers/             # HTTP request handlers
│   ├── services/             # Business logic
│   ├── middleware/           # Auth, logging, CORS
│   └── utils/                # Helper functions
├── pkg/logger/               # Logging package
└── storage/                  # Your uploaded files
```

### Code Flow
```
Request → Middleware (Auth/Logging) → Handler → Service → Storage
                                          ↓
Response ← Middleware ← Handler ← Service ← Storage
```

---

## ❓ Getting Help

- **Issues:** Create an issue on GitHub
- **Questions:** Check API-CONTRACT.md
- **Examples:** See Postman Collection

---

**Happy File Serving! 🎉**
