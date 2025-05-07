# Task 8: Integration Tests

## Alur Kerja

1. Sistem menggunakan dockertest untuk menjalankan tes integrasi
2. Setup tes meliputi:
   - Menjalankan container PostgreSQL dan RabbitMQ menggunakan dockertest
   - Menjalankan migrasi database untuk menyiapkan skema
   - Menginisialisasi aplikasi dengan koneksi ke container tes
   - Menjalankan kasus uji
   - Membersihkan sumber daya setelah tes selesai
3. Tes integrasi mencakup skenario end-to-end untuk memverifikasi fungsionalitas sistem
4. Tes dijalankan secara otomatis sebagai bagian dari pipeline CI/CD

## Endpoint API

Tidak ada endpoint API khusus untuk fitur ini. Tes integrasi menguji endpoint API yang sudah ada.

## Pengujian

1. Jalankan tes integrasi:
   ```bash
   cd backend-golang
   go test -v ./tests/integration/...
   ```

2. Verifikasi bahwa semua kasus uji berhasil, termasuk:
   - Siklus hidup pembuatan dan penghapusan tenant
   - Publikasi dan konsumsi pesan
   - Pembaruan konfigurasi konkurensi
   - Partisi database
   - Cursor pagination

3. Contoh kode setup tes:
   ```go
   pool, err := dockertest.NewPool("")
   resource, _ := pool.Run("postgres", "13", [...])
   
   // Jalankan migrasi
   // Jalankan tes
   defer pool.Purge(resource)
   ```
