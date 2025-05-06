# RabbitMQ Tenant Manager

## Struktur File

### File Utama

- **manager.go**: Berisi definisi struct `TenantManager` dan metode dasar
- **lifecycle.go**: Berisi metode manajemen siklus hidup (`Start`, `Stop`, `StartConsumer`, `StopConsumer`)
- **queue.go**: Berisi metode manajemen queue dan channel
- **consumer_management.go**: Berisi metode manajemen consumer dan fungsi utilitas

### Direktori Consumer

- **consumer/consumer.go**: Berisi pembuatan consumer dan forwarding pesan
- **consumer/worker.go**: Berisi implementasi worker untuk pemrosesan pesan

### Fitur Dead Letter Queue

Implementasi Dead Letter Queue (DLQ) telah dipisahkan ke dalam package terpisah di `pkg/rabbitmq/deadletter.go`, yang menyediakan:

- Retry logic dengan exponential backoff (2, 4, 8 detik)
- Konfigurasi dead letter exchange dan queue
- Konfigurasi TTL pesan (24 jam)
- Konfigurasi jumlah maksimum retry (3 kali)

## Alur Kerja

1. `TenantManager` dibuat menggunakan `NewTenantManager`
2. `Start` dipanggil untuk memulai health check
3. `StartConsumer` dipanggil untuk memulai consumer untuk tenant tertentu
   - Membuat dan mengonfigurasi queue dengan DLQ
   - Membuat worker pool untuk memproses pesan
4. Pesan diterima dari RabbitMQ dan diteruskan ke worker
5. Worker memproses pesan dan menangani error dengan retry logic
6. Jika pemrosesan gagal setelah beberapa kali percobaan, pesan dikirim ke DLQ

## Manajemen Graceful Shutdown

Paket ini terintegrasi dengan `pkg/graceful` untuk mendukung graceful shutdown:

1. `SetShutdownManager` digunakan untuk menetapkan shutdown manager
2. Setiap worker mendaftar ke shutdown manager saat dimulai
3. Setiap worker memberi tahu shutdown manager saat selesai
4. Shutdown manager menunggu semua worker selesai sebelum aplikasi berhenti

## Penggunaan

```go
// Buat tenant manager
tenantManager := rabbitmq.NewTenantManager(rabbitConn, db)

// Set shutdown manager
tenantManager.SetShutdownManager(shutdownManager)

// Mulai tenant manager
tenantManager.Start(ctx)

// Mulai consumer untuk tenant tertentu
tenantManager.StartConsumer(ctx, tenantID)

// Hentikan consumer untuk tenant tertentu
tenantManager.StopConsumer(ctx, tenantID)

// Hentikan tenant manager
tenantManager.Stop(ctx)
```
