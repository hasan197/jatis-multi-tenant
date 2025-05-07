# Task 4: Configurable Transmitter Concurrency

## Alur Kerja

1. Setiap tenant memiliki konfigurasi jumlah worker yang dapat diubah
2. Worker pool diimplementasikan menggunakan pola buffered channel
3. Ketika tenant dibuat, jumlah worker default diambil dari konfigurasi aplikasi
4. Jumlah worker dapat diubah melalui endpoint API khusus
5. Ketika jumlah worker diubah:
   - Sistem menghentikan worker pool yang sedang berjalan
   - Sistem membuat worker pool baru dengan jumlah worker sesuai konfigurasi baru
   - Sistem memulai kembali consumer dengan worker pool baru
6. Variabel atomik digunakan untuk melacak jumlah worker yang aktif

## Endpoint API

- **PUT /api/tenants/{tenant_id}/config/concurrency**
  - **Deskripsi**: Mengubah jumlah worker untuk tenant tertentu
  - **Path Parameter**: 
    - `tenant_id`: UUID tenant yang akan diubah konfigurasinya
  - **Request Body**:
    ```json
    {
      "workers": 5
    }
    ```
  - **Response**: 
    ```json
    {
      "id": "uuid-tenant",
      "name": "Nama Tenant",
      "workers": 5,
      "updated_at": "2023-01-01T00:00:00Z"
    }
    ```

## Pengujian

1. Ubah jumlah worker untuk tenant tertentu:
   ```bash
   curl -X PUT http://localhost:8080/api/tenants/2916830d-8ae9-479f-a5af-5f36cda831de/config/concurrency -H "Content-Type: application/json" -d '{"workers": 5}'
   ```

2. Verifikasi perubahan jumlah worker di database:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT id, name, workers FROM tenants WHERE id = '2916830d-8ae9-479f-a5af-5f36cda831de';"
   ```

3. Verifikasi worker pool telah diperbarui melalui log aplikasi:
   ```bash
   nerdctl logs --tail 50 jatis-sample-stack-golang-backend-golang-1 | grep "\"Started consumer with worker pool\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""
   nerdctl logs --tail 50 jatis-sample-stack-golang-backend-golang-1 | grep "\"Starting worker\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""
   nerdctl logs --tail 50 jatis-sample-stack-golang-backend-golang-1 | grep "Processing message\",\"tenant_id\":\"c383e35c-e199-4ff4-8958-653fbafc9dbc\""
   ```
