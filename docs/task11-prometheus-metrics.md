# Task 11: Prometheus Metrics untuk Monitoring

## Alur Kerja

1. Sistem mengimplementasikan metrics Prometheus untuk monitoring queue depth dan worker activity:
   - Integrasi metrics ke dalam package metrics yang sudah ada di pkg/infrastructure/metrics
   - Menambahkan metrics untuk queue depth, consumer count, worker count, message processing time, retry count, dan dead-lettered messages
   - Memodifikasi worker.go untuk mencatat metrics untuk pemrosesan pesan yang berhasil, gagal, dan retry

2. Metrics diekspos melalui endpoint HTTP khusus:
   - Endpoint /metrics menyediakan semua metrics dalam format Prometheus
   - Prometheus melakukan scraping metrics secara berkala

## Pengujian

1. Jalankan stack monitoring:
   ```bash
   docker-compose -f backend-golang/docker-compose.monitoring.yml up -d

   # atau
   nerdctl compose -f backend-golang/docker-compose.monitoring.yml up -d
   ```

2. Akses Prometheus UI:
   ```
   http://localhost:9090
   ```

3. Verifikasi bahwa metrics diperbarui dengan benar:
   - Kirim beberapa pesan ke RabbitMQ
   - Periksa metrics queue depth di Prometheus
   - Verifikasi bahwa worker activity metrics diperbarui saat pesan diproses

Untuk informasi lebih detail tentang metrics Prometheus, lihat dokumen: [prometheus-metrics-guide.md](../docs/test/prometheus-metrics-guide.md)
