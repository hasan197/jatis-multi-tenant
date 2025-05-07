# Task 2: Auto-Stop Tenant Consumer

## Alur Kerja

1. Ketika permintaan penghapusan tenant diterima melalui endpoint `DELETE /tenants/{id}`:
   - Sistem mengirimkan sinyal shutdown melalui channel kontrol consumer
   - Sistem menutup channel RabbitMQ yang terkait dengan tenant
   - Sistem menghapus antrian RabbitMQ tenant
   - Sistem menghapus tenant dari `TenantManager`
2. Consumer yang sedang berjalan akan menyelesaikan pemrosesan pesan yang sedang ditangani
3. Setelah pemrosesan selesai, consumer akan berhenti dan melepaskan sumber daya
4. Informasi tenant dihapus dari database

## Endpoint API

- **DELETE /api/tenants/{id}**
  - **Deskripsi**: Menghapus tenant dan secara otomatis menghentikan consumer yang terkait
  - **Path Parameter**: 
    - `id`: UUID tenant yang akan dihapus
  - **Response**: 
    ```json
    {
      "message": "Tenant berhasil dihapus"
    }
    ```

## Pengujian

1. Hapus tenant yang ada dengan perintah:
   ```bash
   curl -X DELETE http://localhost:8080/api/tenants/2a7a0324-8118-4e0a-9699-35c1ca694c2e
   ```

2. Verifikasi channel dan consumer telah dihapus dari RabbitMQ:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_channels
   nerdctl exec -it jatis-sample-stack-golang-rabbitmq-1 rabbitmqctl list_consumers
   ```

3. Verifikasi tenant telah dihapus dari database dengan memeriksa tabel tenants
