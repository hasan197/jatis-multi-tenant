# RabbitMQ Publisher UI Documentation

## Overview
Dokumentasi ini menjelaskan tentang UI untuk mempublish message ke RabbitMQ dengan sistem multi-tenant.

## Halaman Utama (Dashboard)

```ascii
+------------------------------------------+
|  RabbitMQ Publisher Dashboard            |
+------------------------------------------+
|                                          |
|  [Pilih Tenant] ▼                        |
|  - Tenant A                              |
|  - Tenant B                              |
|  - Tenant C                              |
|                                          |
|  Status Tenant: [Active]                 |
|  Workers: [3]                            |
|                                          |
+------------------------------------------+
```

## Form Publish Message

```ascii
+------------------------------------------+
|  Publish Message ke Tenant               |
+------------------------------------------+
|                                          |
|  Tenant ID: [2916830d-8ae9-479f-a5af...] |
|                                          |
|  Payload (JSON):                         |
|  +----------------------------------+    |
|  | {                               |    |
|  |   "message": "Hello World",     |    |
|  |   "priority": "high",           |    |
|  |   "metadata": {                 |    |
|  |     "source": "web-ui"          |    |
|  |   }                             |    |
|  | }                               |    |
|  +----------------------------------+    |
|                                          |
|  [Validate JSON]  [Publish Message]      |
|                                          |
+------------------------------------------+
```

## Riwayat Publish

```ascii
+------------------------------------------+
|  Riwayat Publish                         |
+------------------------------------------+
|                                          |
|  Filter:                                 |
|  [Tanggal] [Status] [Tenant]             |
|                                          |
|  +----------------------------------+    |
|  | ID: abc-123                      |    |
|  | Tenant: Tenant A                 |    |
|  | Status: Success                  |    |
|  | Timestamp: 2024-03-20 10:30:00   |    |
|  +----------------------------------+    |
|                                          |
|  [Load More]                            |
|                                          |
+------------------------------------------+
```

## Monitoring

```ascii
+------------------------------------------+
|  Monitoring                              |
+------------------------------------------+
|                                          |
|  Queue Status:                           |
|  - Messages in Queue: 150                |
|  - Active Consumers: 3                   |
|  - Processing Rate: 50 msg/sec           |
|                                          |
|  [Refresh]                               |
|                                          |
+------------------------------------------+
```

## Fitur-fitur Utama

### 1. Pemilihan Tenant
- Dropdown untuk memilih tenant aktif
- Menampilkan status tenant dan jumlah workers
- Validasi tenant status sebelum publish

### 2. Form Publish
- Input JSON payload dengan validasi
- Preview JSON sebelum publish
- Tombol untuk memvalidasi format JSON
- Auto-complete untuk field JSON

### 3. Riwayat Publish
- Tabel dengan pagination
- Filter berdasarkan:
  - Tanggal
  - Status
  - Tenant
- Detail setiap publish message
- Export data ke CSV/Excel

### 4. Monitoring
- Status queue real-time
- Metrik performa
- Jumlah message yang diproses
- Graf visualisasi

## Teknologi yang Digunakan

### Frontend Framework
- React/Vue.js untuk UI yang responsif
- Tailwind CSS untuk styling
- Axios untuk HTTP requests

### Komponen UI
- Code editor untuk JSON (monaco-editor)
- Toast notifications untuk feedback
- Loading states untuk async operations

### State Management
- Redux/Vuex untuk state management
- WebSocket untuk real-time updates

## API Endpoints

```markdown
GET    /api/tenants                    # List semua tenant
POST   /api/messages/{tenant_id}       # Publish message
GET    /api/messages/history           # Riwayat publish
GET    /api/tenants/{id}/status        # Status tenant
PUT    /api/tenants/{id}/config        # Update konfigurasi
```

## Keamanan

### Autentikasi
- JWT token untuk autentikasi
- Role-based access control
- Session management

### Validasi
- Input validation untuk JSON payload
- Rate limiting untuk prevent abuse
- Sanitasi input

### Audit Log
- Logging untuk setiap publish
- Tracking user yang melakukan publish
- Export log untuk audit

## Panduan Penggunaan

1. **Pilih Tenant**
   - Klik dropdown tenant
   - Pilih tenant yang aktif
   - Verifikasi status tenant

2. **Publish Message**
   - Masukkan JSON payload
   - Validasi format JSON
   - Klik tombol publish

3. **Monitor Status**
   - Cek status di halaman monitoring
   - Refresh untuk update real-time
   - Export data jika diperlukan

## Troubleshooting

### Common Issues
1. **Invalid JSON Format**
   - Gunakan JSON validator
   - Periksa sintaks JSON
   - Pastikan semua kurung tutup

2. **Tenant Inactive**
   - Verifikasi status tenant
   - Hubungi admin jika perlu
   - Cek log untuk detail

3. **Publish Failed**
   - Cek koneksi
   - Verifikasi payload
   - Periksa rate limit

## Support

Untuk bantuan teknis, silakan hubungi:
- Email: support@example.com
- Slack: #rabbitmq-support
- Documentation: docs.example.com 