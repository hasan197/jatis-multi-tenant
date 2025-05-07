# Intro: Menjalankan Aplikasi dengan Docker/nerdctl

## Alur Kerja

1. Sistem multi-tenant messaging menggunakan Docker Compose untuk menjalankan semua komponen:
   - Backend Go untuk logika bisnis utama
   - Backend Node.js sebagai proxy API
   - Frontend React untuk antarmuka pengguna
   - PostgreSQL untuk penyimpanan data
   - RabbitMQ untuk messaging
   - Prometheus untuk monitoring
   - CATATAN: developer menggunakan command `nerdctl` sebagai alat docker, untuk command menggunakan `docker` tidak pernah dicoba

2. Siklus hidup aplikasi:
   - Inisialisasi database dan migrasi skema
   - Koneksi ke RabbitMQ
   - Pembuatan exchange dan queue
   - Menjalankan consumer untuk tenant yang aktif
   - Menyediakan API untuk interaksi dengan sistem

## Endpoint API

Setelah aplikasi berjalan, endpoint API berikut tersedia:

- **Backend Go**: `http://localhost:8080`
  - Endpoint API utama untuk manajemen tenant dan pesan
  - Endpoint metrics Prometheus: `http://localhost:8080/metrics`
  - Dokumentasi Swagger: `http://localhost:8080/swagger/index.html`

- **Frontend React**: `http://localhost:5173`
  - Antarmuka pengguna untuk interaksi dengan sistem
  - Form untuk mempublikasikan pesan

- **Proxy Server**: `http://localhost:3000`
  - API untuk interaksi dengan sistem dan diteruskan ke Backend Go
  - API untuk publish pesan ke RabbitMQ http://localhost:3000/api/tenants/{TENANT_ID}/publish

## Pengujian

### Menjalankan Aplikasi

1. Menjalankan aplikasi lengkap dengan Docker Compose:
   ```bash
   # Menggunakan Docker
   docker-compose up -d --build 
   
   # atau tanpa --build jika tanpa rebuild
   
   docker-compose up -d

   # Menggunakan nerdctl
   nerdctl compose up -d --build
   
   # atau tanpa --build jika tanpa rebuild
   
   nerdctl compose up -d
   ```

2. Menjalankan stack monitoring:
   ```bash
   # Menggunakan Docker
   docker-compose -f backend-golang/docker-compose.monitoring.yml up -d

   # Menggunakan nerdctl
   nerdctl compose -f backend-golang/docker-compose.monitoring.yml up -d
   ```

### Memeriksa Status Layanan

1. Memeriksa status container:
   ```bash
   # Menggunakan Docker
   docker-compose ps

   # Menggunakan nerdctl
   nerdctl compose ps
   ```

2. Melihat log aplikasi:
   ```bash
   # Menggunakan Docker
   docker-compose logs -f backend-golang

   # Menggunakan nerdctl
   nerdctl logs -f jatis-sample-stack-golang-backend-golang-1
   ```

3. Memeriksa database:
   ```bash
   # Menggunakan Docker
   docker exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db

   # Menggunakan nerdctl
   nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db
   ```

4. Memeriksa RabbitMQ:
   ```bash
   # Menggunakan Docker
   docker exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_queues

   # Menggunakan nerdctl
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_queues
   ```

### Menghentikan Aplikasi

1. Menghentikan dan menghapus container:
   ```bash
   # Menggunakan Docker
   docker-compose down

   # Menggunakan nerdctl
   nerdctl compose down
   ```

2. Menghentikan, menghapus container, dan menghapus volume:
   ```bash
   # Menggunakan Docker
   docker-compose down -v

   # Menggunakan nerdctl
   nerdctl compose down -v
   ```


### Dokumentasi Tugas
Cek folder [docs/](docs/)