# Panduan Monitoring dengan Prometheus Metrics

Dokumen ini menjelaskan cara menggunakan dan memantau metrics Prometheus untuk sistem multi-tenant messaging dengan RabbitMQ.

## Prasyarat

- Docker dan Docker Compose
- Aplikasi backend-golang, backend-nodejs, dan RabbitMQ yang berjalan
- Prometheus dan Grafana yang dikonfigurasi

## Metrics yang Tersedia

### HTTP Metrics

| Nama | Tipe | Deskripsi | Label |
|------|------|-----------|-------|
| `http_requests_total` | Counter | Total jumlah request HTTP | `method`, `endpoint`, `status` |
| `http_request_duration_seconds` | Histogram | Durasi request HTTP dalam detik | `method`, `endpoint` |

### RabbitMQ Queue Metrics

| Nama | Tipe | Deskripsi | Label |
|------|------|-----------|-------|
| `rabbitmq_queue_depth` | Gauge | Jumlah pesan dalam queue | `tenant_id`, `queue_name` |
| `rabbitmq_queue_consumer_count` | Gauge | Jumlah consumer untuk queue | `tenant_id`, `queue_name` |

### RabbitMQ Worker Metrics

| Nama | Tipe | Deskripsi | Label |
|------|------|-----------|-------|
| `rabbitmq_worker_count` | Gauge | Jumlah worker aktif | `tenant_id` |
| `rabbitmq_messages_processed_total` | Counter | Total jumlah pesan yang diproses | `tenant_id`, `status` |
| `rabbitmq_message_processing_time_seconds` | Histogram | Waktu pemrosesan pesan dalam detik | `tenant_id` |

### Dead Letter Queue (DLQ) Metrics

| Nama | Tipe | Deskripsi | Label |
|------|------|-----------|-------|
| `rabbitmq_dlq_depth` | Gauge | Jumlah pesan dalam dead letter queue | `tenant_id` |
| `rabbitmq_message_retry_total` | Counter | Total jumlah retry pesan | `tenant_id` |
| `rabbitmq_messages_dead_lettered_total` | Counter | Total jumlah pesan yang dikirim ke DLQ | `tenant_id` |

### Database Metrics

| Nama | Tipe | Deskripsi | Label |
|------|------|-----------|-------|
| `db_query_duration_seconds` | Histogram | Durasi query database dalam detik | `query_type` |

## Cara Menjalankan Prometheus dan Grafana

Prometheus dan Grafana sudah dikonfigurasi dalam file `docker-compose.monitoring.yml`. Untuk menjalankannya:

```bash
# Dari direktori root proyek
cd backend-golang
docker-compose -f docker-compose.monitoring.yml up -d
```

Setelah container berjalan, Anda dapat mengakses:
- Prometheus UI: http://localhost:9090
- Grafana UI: http://localhost:3000 (username: admin, password: admin123)

## Endpoint Metrics

Aplikasi backend-golang mengekspos metrics Prometheus di endpoint:

```
http://localhost:8080/metrics
```

Prometheus dikonfigurasi untuk scrape metrics dari endpoint ini setiap 15 detik.

## Melihat Metrics di Prometheus

1. Buka Prometheus UI di http://localhost:9090
2. Gunakan tab "Graph" untuk melihat metrics
3. Ketik nama metrics di field "Expression" (misalnya `rabbitmq_queue_depth`)
4. Klik tombol "Execute" untuk melihat data

## Dashboard Grafana

Grafana sudah dikonfigurasi dengan datasource Prometheus. Untuk membuat dashboard:

1. Buka Grafana UI di http://localhost:3000
2. Login dengan username `admin` dan password `admin123`
3. Klik "Create Dashboard"
4. Tambahkan panel baru dan pilih metrics yang ingin ditampilkan

## Contoh Query PromQL

Berikut adalah beberapa contoh query PromQL yang berguna:

### Melihat Queue Depth per Tenant

```
rabbitmq_queue_depth{tenant_id="$tenant_id"}
```

### Melihat Rate Pemrosesan Pesan

```
rate(rabbitmq_messages_processed_total{status="success"}[5m])
```

### Melihat Jumlah Pesan di DLQ

```
rabbitmq_dlq_depth
```

### Melihat Waktu Pemrosesan Pesan (99th percentile)

```
histogram_quantile(0.99, sum(rate(rabbitmq_message_processing_time_seconds_bucket[5m])) by (le, tenant_id))
```

## Alerting

Anda dapat mengonfigurasi alert di Prometheus atau Grafana untuk memantau kondisi tertentu, seperti:

- Queue depth yang terlalu tinggi
- Pesan di DLQ yang bertambah
- Waktu pemrosesan pesan yang terlalu lama
- Rate error yang tinggi

## Troubleshooting

### Metrics Tidak Muncul di Prometheus

1. Pastikan endpoint `/metrics` dapat diakses
2. Periksa konfigurasi Prometheus di `prometheus.yml`
3. Periksa log Prometheus untuk error

### Grafana Tidak Menampilkan Data

1. Pastikan datasource Prometheus dikonfigurasi dengan benar
2. Periksa query PromQL yang digunakan
3. Pastikan Prometheus memiliki data untuk metrics yang diminta

## Kesimpulan

Dengan menggunakan Prometheus dan Grafana, Anda dapat memantau performa sistem messaging multi-tenant dengan RabbitMQ secara real-time. Metrics yang disediakan mencakup queue depth, worker activity, dan DLQ status, yang memungkinkan Anda untuk mengidentifikasi bottleneck dan masalah potensial sebelum mempengaruhi pengguna.
