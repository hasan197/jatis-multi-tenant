# Panduan Pengujian Graceful Shutdown

Dokumen ini menjelaskan cara menguji fitur graceful shutdown pada aplikasi backend-golang dengan RabbitMQ.

## Prasyarat

- Docker/nerdctl untuk menjalankan container
- curl untuk mengirim request HTTP
- Terminal/command prompt

## Langkah-langkah Pengujian

### 1. Memastikan Aplikasi Berjalan

```bash
# Periksa status container
nerdctl ps | grep backend-golang

# Jika container tidak berjalan, jalankan dengan docker-compose
nerdctl compose up -d
```

### 2. Membuat Tenant Baru untuk Pengujian

```bash
# Buat tenant baru untuk pengujian
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Graceful Shutdown Test","description":"Tenant untuk menguji fitur graceful shutdown","status":"active","workers":5}' \
  http://localhost:8080/api/tenants

# Catat ID tenant yang dikembalikan untuk digunakan pada langkah selanjutnya
# Contoh output: {"id":"c383e35c-e199-4ff4-8958-653fbafc9dbc",...}
```

### 3. Memeriksa Status Tenant dan Worker

```bash
# Periksa detail tenant
curl -s http://localhost:8080/api/tenants/TENANT_ID | jq

# Periksa status antrian tenant
curl -s http://localhost:8080/api/tenants/TENANT_ID/queue-status | jq

# Periksa consumer yang aktif
curl -s http://localhost:8080/api/tenants/consumers | jq
```

### 4. Mengirim Pesan yang Membutuhkan Waktu untuk Diproses

```bash
# Kirim 10 pesan dengan waktu pemrosesan 10 detik
for i in {1..10}; do 
  curl -X POST -H "Content-Type: application/json" \
    -d "{\"message\": \"Long processing message $i\", \"priority\": \"high\", \"processing_time\": 10}" \
    http://localhost:3000/api/tenants/TENANT_ID/publish
  echo
done
```

### 5. Memeriksa Log untuk Memastikan Pesan Sedang Diproses

```bash
# Periksa log untuk melihat pesan yang sedang diproses
nerdctl logs --tail 20 jatis-sample-stack-golang-backend-golang-1
```

Anda akan melihat log seperti:
```
{"level":"debug","message_id":"","msg":"Processing message","tenant_id":"TENANT_ID","time":"2025-05-06T10:58:00.197Z","worker_id":0}
{"level":"debug","message_id":"","msg":"Processing message","tenant_id":"TENANT_ID","time":"2025-05-06T10:58:00.255Z","worker_id":2}
...
```

### 6. Mengirim Sinyal SIGTERM untuk Memicu Graceful Shutdown

```bash
# Kirim sinyal SIGTERM ke container backend-golang
nerdctl kill --signal SIGTERM jatis-sample-stack-golang-backend-golang-1
```

### 7. Memeriksa Log untuk Memverifikasi Proses Graceful Shutdown

```bash
# Periksa log untuk melihat proses graceful shutdown
nerdctl logs jatis-sample-stack-golang-backend-golang-1 | grep -A 20 "Sinyal shutdown diterima"
```

Anda akan melihat log seperti:
```
{"level":"info","msg":"Sinyal shutdown diterima, memulai proses graceful shutdown","time":"2025-05-06T10:58:28.160Z"}
{"level":"info","msg":"Menghentikan HTTP server (tidak menerima request baru)","time":"2025-05-06T10:58:28.160Z"}
{"level":"info","msg":"Menunggu semua proses aktif selesai","time":"2025-05-06T10:58:28.162Z"}
{"level":"info","msg":"Semua proses aktif selesai dengan sukses","time":"2025-05-06T10:58:28.162Z"}
{"level":"info","msg":"Menutup semua resource","time":"2025-05-06T10:58:28.162Z"}
...
{"level":"info","msg":"Graceful shutdown selesai","time":"2025-05-06T10:58:28.180Z"}
Server shutdown complete
```

## Verifikasi Hasil

Proses graceful shutdown berjalan dengan benar jika:

1. Server berhenti menerima request baru setelah menerima sinyal SIGTERM
2. Server menunggu semua proses yang sedang berjalan untuk selesai
3. Server menutup semua resource dengan rapi (RabbitMQ channel, worker, dll)
4. Server keluar dengan status sukses

## Catatan Penting

- Pastikan untuk mengganti `TENANT_ID` dengan ID tenant yang sebenarnya dari hasil langkah 2
- Waktu pemrosesan pesan dapat diatur melalui parameter `processing_time` (dalam detik)
- Jumlah worker dapat diatur melalui parameter `workers` saat membuat tenant atau dengan endpoint konfigurasi concurrency:

```bash
# Mengubah jumlah worker untuk tenant
curl -X PUT http://localhost:8080/api/tenants/TENANT_ID/config/concurrency \
  -H "Content-Type: application/json" \
  -d '{"workers": 5}'
```
