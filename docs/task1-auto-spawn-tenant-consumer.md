# Task 1: Auto-Spawn Tenant Consumer

## Alur Kerja

1. Sistem menggunakan `TenantManager` untuk melacak tenant aktif dan consumer RabbitMQ mereka.
2. Ketika tenant baru dibuat melalui endpoint `POST /tenants`:
   - Sistem membuat antrian RabbitMQ khusus dengan format `tenant_{id}_queue`
   - Sistem menjalankan goroutine consumer yang mendengarkan antrian tersebut
   - Sistem menyimpan channel kontrol consumer di `TenantManager`
3. Setiap consumer menggunakan `channel.Consume()` dengan tag consumer unik
4. Pesan yang diterima dari antrian diproses oleh worker pool
5. Setiap pesan yang berhasil diproses akan disimpan ke dalam tabel partisi database khusus tenant

## Endpoint API

- **POST /api/tenants**
  - **Deskripsi**: Membuat tenant baru dan secara otomatis menjalankan consumer untuk tenant tersebut
  - **Request Body**:
    ```json
    {
      "name": "Nama Tenant",
      "description": "Deskripsi tenant",
      "status": "active"
    }
    ```
  - **Response**: 
    ```json
    {
      "id": "uuid-tenant",
      "name": "Nama Tenant",
      "description": "Deskripsi tenant",
      "status": "active",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
    ```

## Pengujian

1. Buat tenant baru dengan perintah:
   ```bash
   curl -X POST http://localhost:8080/api/tenants -H "Content-Type: application/json" -d '{"name":"Test Partition","description":"Test tenant for partition verification","status":"active"}'
   ```

2. Verifikasi tabel partisi telah dibuat:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT tablename FROM pg_tables WHERE tablename LIKE 'messages_%';"
    # atau
   docker exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT tablename FROM pg_tables WHERE tablename LIKE 'messages_%';"
   ```

3. Verifikasi consumer telah berjalan di RabbitMQ:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_channels
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_consumers

   # atau
   docker exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_channels
   docker exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_consumers
   ```
