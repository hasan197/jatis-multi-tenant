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