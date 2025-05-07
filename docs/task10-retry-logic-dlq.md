# Task 10: Retry Logic dan Dead-Letter Queue

## Alur Kerja

1. Sistem mengimplementasikan mekanisme retry dan Dead-Letter Queue (DLQ) untuk menangani pesan yang gagal diproses:
   - Membuat dead-letter exchange (DLX) bernama "dlx.tenant"
   - Membuat dead-letter queue untuk setiap tenant dengan format "dlq.tenant.{tenant_id}"
   - Mengonfigurasi main queue dengan parameter x-dead-letter-exchange dan x-dead-letter-routing-key
   - Menambahkan TTL 24 jam untuk pesan di queue

2. Mekanisme retry bekerja dengan cara:
   - Menambahkan header "x-retry-count" untuk melacak jumlah percobaan
   - Implementasi maksimal 3 kali percobaan (maxRetries)
   - Menggunakan exponential backoff (2, 4, 8 detik) untuk delay antar percobaan
   - Mempublikasi ulang pesan dengan retry count yang diperbarui

3. Penanganan kegagalan:
   - Jika masih dalam batas retry, pesan akan dicoba ulang dengan delay
   - Jika sudah mencapai batas retry, pesan akan dikirim ke dead-letter queue
   - Pesan di dead-letter queue dapat diproses secara manual atau dengan consumer terpisah

## Endpoint API

Tidak ada endpoint API khusus untuk fitur ini. Retry logic dan DLQ diimplementasikan sebagai bagian dari infrastruktur messaging.

## Pengujian

1. Jalankan skrip pengujian DLQ:
   ```bash
   ./test-dlq.sh
   ```

2. Periksa status dead-letter queue:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_queues | grep dlq

   #atau
   docker exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_queues | grep dlq
   ```

Untuk informasi lebih detail tentang pengujian DLQ, lihat dokumen: [dlq-testing-guide.md](../docs/test/dlq-testing-guide.md)
