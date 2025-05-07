# Panduan Pengujian Dead Letter Queue (DLQ)

Dokumen ini menjelaskan cara menguji fitur Dead Letter Queue (DLQ) pada aplikasi multi-tenant dengan RabbitMQ.

## Prasyarat

- curl untuk mengirim request HTTP
- Terminal/command prompt
- Aplikasi backend-golang dan backend-nodejs yang berjalan
- RabbitMQ yang berjalan

## Konsep DLQ

Dead Letter Queue (DLQ) adalah mekanisme untuk menangani pesan yang gagal diproses setelah beberapa kali percobaan. Alur kerjanya:

1. Pesan dikirim ke queue reguler
2. Consumer mencoba memproses pesan
3. Jika gagal, pesan akan dicoba ulang beberapa kali (retry)
4. Jika tetap gagal setelah batas retry tercapai, pesan akan dikirim ke DLQ
5. Pesan di DLQ dapat diproses secara manual atau dengan consumer terpisah

## Konfigurasi DLQ

Aplikasi ini menggunakan konfigurasi DLQ sebagai berikut:

- Dead Letter Exchange (DLX): `dlx.tenant`
- Dead Letter Queue: `dlq.tenant.{tenant_id}`
- Max Retries: 3
- Retry Delay: Exponential backoff (2, 4, 8 detik)
- TTL untuk pesan di queue: 24 jam (86400000 ms)

## Langkah-langkah Pengujian

### 1. Memastikan Aplikasi Berjalan

```bash
# Periksa status container
nerdctl ps | grep jatis-sample-stack-golang

# Jika container tidak berjalan, jalankan dengan docker-compose
nerdctl compose up -d
```

### 2. Membuat Tenant Baru untuk Pengujian

```bash
# Buat tenant baru untuk pengujian
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Test Tenant DLQ","description":""}' \
  http://localhost:8080/api/tenants

# Catat ID tenant yang dikembalikan untuk digunakan pada langkah selanjutnya
# Contoh output: {"id":"a22ca499-8f02-4dc1-a297-d18bbe2e4ac7",...}
```

### 3. Mempublikasikan Pesan dengan Flag Force Error

```bash
# Publikasikan pesan dengan flag force_error untuk memicu error
curl -X POST -H "Content-Type: application/json" \
  -d '{"content": "Test message with error", "metadata": {"force_error": true}}' \
  http://localhost:3000/api/tenants/TENANT_ID/publish

# Ganti TENANT_ID dengan ID tenant yang dibuat pada langkah sebelumnya
```

### 4. Memantau Proses Retry dan DLQ

```bash
# Lihat log dari backend-golang untuk memantau proses retry dan DLQ
nerdctl logs jatis-sample-stack-golang-backend-golang-1 | grep -i "DLQ\|retry\|force_error"
```

Anda akan melihat log yang menunjukkan:
- Pesan diterima dengan flag `force_error`
- Error simulasi dipicu
- Pesan di-NACK dan di-requeue untuk retry
- Setelah batas retry tercapai, pesan dikirim ke DLQ

### 5. Memeriksa Status Queue

```bash
# Periksa status queue untuk tenant
curl http://localhost:8080/api/tenants/TENANT_ID/queue-status

# Ganti TENANT_ID dengan ID tenant yang dibuat pada langkah sebelumnya
```

### 6. Menggunakan RabbitMQ Management UI

Anda juga dapat menggunakan RabbitMQ Management UI untuk memeriksa queue dan pesan:

1. Buka http://localhost:15672/ di browser
2. Login dengan username `guest` dan password `guest`
3. Navigasi ke tab "Queues"
4. Periksa queue `tenant.{tenant_id}` dan `dlq.tenant.{tenant_id}`

## Menggunakan Script Pengujian Otomatis

Untuk memudahkan pengujian, Anda dapat menggunakan script `test-dlq.sh` yang disediakan:

```bash
# Jalankan script pengujian
./test-dlq.sh
```

Script ini akan:
1. Membuat tenant baru
2. Mempublikasikan pesan dengan flag `force_error`
3. Menunggu proses retry
4. Memeriksa status queue
5. Memeriksa log untuk memverifikasi alur DLQ

## Troubleshooting

### Pesan Tidak Masuk ke DLQ

1. Pastikan konfigurasi DLQ di backend Go sudah benar
2. Periksa apakah flag `force_error` terdeteksi dengan benar
3. Periksa apakah mekanisme NACK berfungsi dengan benar
4. Periksa log untuk error atau warning

### Error Publikasi Pesan

1. Pastikan backend Node.js tidak mencoba mendeklarasikan queue yang sudah ada
2. Pastikan format pesan sesuai dengan yang diharapkan oleh worker Go

## Implementasi DLQ

Implementasi DLQ terdapat di beberapa file:

- `pkg/rabbitmq/deadletter.go`: Implementasi utama DLQ
- `internal/modules/tenant/delivery/messaging/rabbitmq/consumer/worker.go`: Implementasi worker yang mendeteksi flag `force_error`
- `backend-nodejs/src/index.ts`: Implementasi publikasi pesan dari Node.js

## Kesimpulan

Dengan mengikuti langkah-langkah di atas, Anda dapat menguji dan memverifikasi bahwa mekanisme Dead Letter Queue berfungsi dengan baik. Pesan yang gagal diproses akan dicoba ulang beberapa kali sebelum akhirnya dikirim ke DLQ, memastikan tidak ada pesan yang hilang dan memberikan kesempatan untuk menangani kegagalan pemrosesan dengan lebih baik.
