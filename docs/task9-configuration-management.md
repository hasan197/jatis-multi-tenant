# Task 9: Configuration Management

## Alur Kerja

1. Sistem menggunakan file konfigurasi YAML untuk mengatur parameter aplikasi (backend-golang/configs/config.yaml)
2. Konfigurasi mencakup:
   - Koneksi RabbitMQ
   - Koneksi database
   - Jumlah worker default
   - Parameter aplikasi lainnya
3. Konfigurasi dimuat saat aplikasi dimulai

## Endpoint API

Tidak ada endpoint API khusus untuk fitur ini. Konfigurasi dikelola melalui file konfigurasi.

Struktur file konfigurasi:
```yaml
app:
  name: sample-stack-golang
  port: 8080
  version: 1.0.0
  env: development
  workers: 2 # Default worker count
```

## Pengujian

1. Modifikasi file konfigurasi:
   ```bash
   nano backend-golang/configs/config.yaml
   ```

2. Ubah nilai parameter, misalnya jumlah worker default:
   ```yaml
   app:
     workers: 5
   ```

3. Restart aplikasi:
   ```bash
   docker-compose restart backend-golang
   ```

4. Verifikasi bahwa aplikasi menggunakan nilai konfigurasi baru:
   - Buat tenant baru dan verifikasi jumlah worker sesuai dengan nilai default baru
   - Periksa log aplikasi untuk memastikan konfigurasi dimuat dengan benar
      ```
      nerdctl logs --tail 50 jatis-sample-stack-golang-backend-golang-1 |grep "\"Started consumer with worker pool\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""

      nerdctl logs --tail 50 jatis-sample-stack-golang-backend-golang-1 |grep "\"Starting worker\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""

      # atau

      docker logs --tail 50 jatis-sample-stack-golang-backend-golang-1 |grep "\"Started consumer with worker pool\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""

      docker logs --tail 50 jatis-sample-stack-golang-backend-golang-1 |grep "\"Starting worker\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""

      ```