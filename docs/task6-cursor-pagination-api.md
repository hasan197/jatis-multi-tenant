# Task 6: Cursor Pagination API

## Alur Kerja

1. Sistem mengimplementasikan cursor pagination untuk mengambil data pesan
2. Ketika endpoint `/api/messages` dipanggil:
   - Sistem mengambil parameter `cursor` dan `limit` dari query string
   - Jika `cursor` tidak disediakan, sistem mengambil data dari awal
   - Jika `cursor` disediakan, sistem mengambil data setelah cursor tersebut
   - Sistem membatasi jumlah data yang diambil sesuai parameter `limit`
   - Sistem mengembalikan data pesan beserta cursor berikutnya
3. Cursor diimplementasikan menggunakan ID pesan terakhir yang diambil
4. Implementasi dilakukan di tiga layer:
   - Handler di message_handler.go dengan method GetMessages
   - Repository di message_repository.go dengan method FindAll
   - UseCase di message_usecase.go dengan method GetMessages

## Endpoint API

- **GET /api/messages**
  - **Deskripsi**: Mengambil daftar pesan dengan pagination berbasis cursor
  - **Query Parameters**: 
    - `cursor`: Cursor untuk pagination (opsional)
    - `limit`: Jumlah maksimum data yang diambil (opsional, default: 10)
  - **Response**: 
    ```json
    {
      "data": [
        {
          "id": "uuid-pesan-1",
          "tenant_id": "uuid-tenant",
          "payload": { "key": "value" },
          "created_at": "2023-01-01T00:00:00Z"
        },
        {
          "id": "uuid-pesan-2",
          "tenant_id": "uuid-tenant",
          "payload": { "key": "value" },
          "created_at": "2023-01-01T00:00:00Z"
        }
      ],
      "next_cursor": "uuid-pesan-2"
    }
    ```

## Pengujian

1. Ambil halaman pertama pesan:
   ```bash
   curl -X GET "http://localhost:8080/api/messages?limit=5"
   ```

2. Ambil halaman berikutnya menggunakan cursor dari respons sebelumnya:
   ```bash
   curl -X GET "http://localhost:8080/api/messages?cursor=uuid-pesan-terakhir&limit=5"
   ```

3. Verifikasi bahwa:
   - Respons berisi jumlah pesan sesuai parameter `limit`
   - Pesan diurutkan berdasarkan waktu pembuatan
   - Cursor berikutnya menunjuk ke ID pesan terakhir
   - Tidak ada duplikasi pesan antar halaman

Untuk informasi lebih detail tentang pengujian cursor pagination, lihat dokumen: [cursor-pagination-test.md](../docs/test/cursor-pagination-test.md)
