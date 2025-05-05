# Jatis Sample Stack Golang

## Cara Menjalankan Aplikasi

1. Jalankan stack aplikasi utama (backend, frontend, database, dsb):
   ```bash
   nerdctl compose up -d
   # atau
   docker compose up -d
   ```

2. Untuk menghentikan stack aplikasi:
   ```bash
   nerdctl compose down
   # atau
   docker compose down
   ```

## Cara Menjalankan Monitoring (Prometheus & Grafana)

1. Jalankan stack monitoring:
   ```bash
   nerdctl compose -f backend-golang/docker-compose.monitoring.yml up -d
   # atau
   docker compose -f backend-golang/docker-compose.monitoring.yml up -d
   ```

2. Untuk menghentikan stack monitoring:
   ```bash
   nerdctl compose -f backend-golang/docker-compose.monitoring.yml down
   # atau
   docker compose -f backend-golang/docker-compose.monitoring.yml down
   ```

## Akses Monitoring

- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3000 (user: admin, password: admin123)
- **Alertmanager:** http://localhost:9093

> **Catatan:**
> Stack aplikasi dan stack monitoring dipisahkan agar lebih modular dan mudah di-maintain. Jalankan keduanya sesuai kebutuhan.

---

## Cara Menjalankan Test Backend Golang

1. Build image backend-golang:
   ```bash
   nerdctl compose -f backend-golang/docker-compose.test.yml build
   ```

2. Jalankan test backend-golang:
   ```bash
   nerdctl compose -f backend-golang/docker-compose.test.yml up --abort-on-container-exit
   # atau
   docker compose -f backend-golang/docker-compose.test.yml up --abort-on-container-exit
   ```

3. Hasil test akan muncul di terminal. Untuk menghentikan dan membersihkan container test:
   ```bash
   nerdctl compose -f backend-golang/docker-compose.test.yml down
   # atau
   docker compose -f backend-golang/docker-compose.test.yml down
   ```

## Cara Update Dependencies Go

Untuk memperbarui dependencies Go tanpa perlu menginstall Go di komputer lokal, gunakan perintah berikut:

```bash
nerdctl run --rm -v $(pwd)/backend-golang:/app -w /app golang:1.21-alpine sh -c "apk add --no-cache git && go mod tidy"
```

## Cara Menjalankan Perintah di Container Menggunakan nerdctl

### Menjalankan Query SQL

Untuk menjalankan query SQL di database PostgreSQL:

```bash
# Menjalankan query SQL langsung
nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT * FROM users;"
nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "\dt;"

# Menjalankan file SQL
nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db < path/to/query.sql
```

### Menjalankan Perintah Go

Untuk menjalankan perintah Go di container:

```bash
# Menjalankan go mod tidy
nerdctl exec -it jatis-sample-stack-golang-backend-golang-1 go mod tidy

# Menjalankan go test
nerdctl exec -it jatis-sample-stack-golang-backend-golang-1go test ./...

# Menjalankan go build
nerdctl exec -it jatis-sample-stack-golang-backend-golang-1 go build
```

### Menjalankan Perintah di Container Lainnya

```bash
# Menjalankan perintah di container Redis
nerdctl exec -it redis redis-cli

# Menjalankan perintah di container Prometheus
nerdctl exec -it prometheus promtool check config /etc/prometheus/prometheus.yml

# Menjalankan perintah di container Grafana
nerdctl exec -it grafana grafana-cli plugins list
```

### Sistem Messaging (RabbitMQ)

Sistem messaging menggunakan RabbitMQ untuk menangani komunikasi antar tenant. Setiap tenant memiliki queue dan consumer sendiri.

#### Struktur Queue
- Format nama queue: `tenant.{tenant_id}`
- Exchange: `amq.default`
- Routing key: `tenant.{tenant_id}`

#### Perintah RabbitMQ

```bash
# Melihat daftar queue
nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_queues

# Melihat daftar consumer
nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_consumers

# Mengirim pesan ke queue tenant
nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqadmin publish exchange=amq.default routing_key="tenant.{tenant_id}" payload='{"text": "Pesan test"}'

# Melihat status consumer tenant
curl http://localhost:8080/api/tenants/{tenant_id}/consumers
```

#### Fitur Consumer
1. **Auto Creation**: Consumer dibuat otomatis saat tenant dibuat
2. **Health Check**: Status consumer dicek setiap 30 detik
3. **Auto Recovery**: Consumer yang tidak aktif akan di-restart otomatis
4. **Cleanup**: Consumer dihapus saat tenant dihapus

#### Format Pesan
```json
{
    "text": "Isi pesan"
}
```

> **Catatan:**
> - Ganti `backend-golang`, `postgres`, `redis`, dll dengan nama container yang sesuai
> - Gunakan flag `-it` jika perintah membutuhkan interaksi
> - Gunakan flag `-i` jika perintah hanya membutuhkan input
> - Gunakan flag `-v` untuk mount volume jika diperlukan

