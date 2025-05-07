# Task 5: Graceful Shutdown

## Alur Kerja

1. Ketika aplikasi menerima sinyal shutdown (SIGINT atau SIGTERM):
   - Sistem menghentikan penerimaan permintaan HTTP baru
   - Sistem mengirimkan sinyal shutdown ke semua consumer yang aktif
   - Sistem menunggu semua transaksi yang sedang berjalan selesai diproses
   - Sistem menutup semua koneksi database dan RabbitMQ
   - Sistem keluar dengan status sukses
2. Proses shutdown menggunakan context dengan timeout untuk memastikan aplikasi tidak menunggu terlalu lama
3. Setiap consumer yang menerima sinyal shutdown:
   - Menyelesaikan pemrosesan pesan yang sedang ditangani
   - Menutup channel RabbitMQ
   - Menghentikan worker pool
   - Mengirimkan sinyal selesai ke aplikasi utama

## Endpoint API

Tidak ada endpoint API khusus untuk fitur ini. Graceful shutdown diimplementasikan sebagai bagian dari siklus hidup aplikasi dan dipicu oleh sinyal sistem operasi.

## Pengujian

1. Jalankan aplikasi dengan Docker Compose:
   ```bash
   nerdctl compose up -d
   #atau
   docker-compose up -d
   ```

2. Kirim beberapa pesan ke antrian RabbitMQ untuk diproses

3. Kirim sinyal SIGTERM ke container aplikasi:
   ```bash
   nerdctl kill --signal=SIGTERM jatis-sample-stack-golang-backend-golang-1
   #atau
   docker kill --signal=SIGTERM jatis-sample-stack-golang-backend-golang-1
   ```

4. Verifikasi melalui log bahwa aplikasi:
   - Menerima sinyal shutdown
   - Menyelesaikan pemrosesan pesan yang sedang berjalan
   - Menutup semua koneksi dengan baik
   - Keluar dengan status sukses

5. Verifikasi tidak ada pesan yang hilang atau rusak di RabbitMQ

Untuk informasi lebih detail tentang pengujian graceful shutdown, lihat dokumen: [graceful-shutdown-test.md](../docs/test/graceful-shutdown-test.md)
