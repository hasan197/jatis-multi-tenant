# Panduan Docker Compose

## Overview

Proyek ini menggunakan Docker dan Docker Compose untuk mengelola lingkungan pengembangan dan produksi. Semua layanan dikonfigurasi dengan multi-stage build untuk mendukung development dan production environment.

## Struktur Docker

Semua layanan menggunakan pendekatan multi-stage build dengan tiga tahap utama:
1. **Builder**: Tahap untuk membangun aplikasi
2. **Development**: Tahap untuk pengembangan dengan hot-reload
3. **Production**: Tahap untuk produksi dengan optimasi

## Cara Penggunaan

### Development Environment

Untuk menjalankan semua layanan dalam mode development:

```bash
nerdctl compose up -d
```

Atau untuk menjalankan layanan tertentu:

```bash
nerdctl compose up -d backend-golang
```

### Production Environment

Untuk menjalankan dalam mode production:

```bash
nerdctl compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Rebuild Image

Untuk membangun ulang image setelah perubahan Dockerfile:

```bash
nerdctl compose build
```

Atau untuk layanan tertentu:

```bash
nerdctl compose build backend-golang
```

### Logs

Untuk melihat logs:

```bash
nerdctl compose logs
```

Atau untuk layanan tertentu:

```bash
nerdctl compose logs -f backend-golang
```

## Service Endpoints

- **Backend Golang**: http://localhost:8080
- **Backend NodeJS**: http://localhost:3000
- **Frontend React**: http://localhost:5173
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **RabbitMQ**: localhost:5672 (management: http://localhost:15672) 