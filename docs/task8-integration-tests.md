# Task 8: Integration Tests

## Alur Kerja

1. Sistem menggunakan dockertest untuk menjalankan tes integrasi
2. Setup tes meliputi:
   - Menjalankan container PostgreSQL, Redis, dan RabbitMQ menggunakan dockertest
   - Menjalankan migrasi database untuk menyiapkan skema
   - Menginisialisasi aplikasi dengan koneksi ke container tes
   - Menjalankan test case
   - Membersihkan sumber daya setelah tes selesai
3. Tes integrasi mencakup skenario end-to-end untuk memverifikasi fungsionalitas sistem

## Struktur Tes Integrasi

```
backend-golang/test/integration/
├── go.mod                 # Modul Go untuk tes integrasi
├── go.sum                 # Dependensi untuk tes integrasi
├── message_test.go        # Tes untuk modul pesan
├── setup/                 # Konfigurasi pengujian
│   ├── docker.go          # Setup container Docker untuk tes
│   └── logger.go          # Konfigurasi logger untuk tes
└── tenant_test.go         # Tes untuk modul tenant
```

## Fitur yang Diuji

1. **Tenant Management**:
   - Pembuatan dan penghapusan tenant
   - Manajemen konfigurasi tenant

2. **Message Management**:
   - Pembuatan pesan dalam partisi tenant
   - Pengambilan pesan berdasarkan tenant
   - Partisi database per tenant

3. **RabbitMQ Integration**:
   - Koneksi ke RabbitMQ
   - Pembuatan antrian berdasarkan tenant

## Pengujian

1. Jalankan tes integrasi:
   ```bash
   cd backend-golang/test/integration/
   go test -v ./...
   ```

2. Untuk menjalankan dengan Rancher Desktop, gunakan:
   ```bash
   DOCKER_HOST=unix:///Users/admin/.rd/docker.sock go test -v ./...
   ```