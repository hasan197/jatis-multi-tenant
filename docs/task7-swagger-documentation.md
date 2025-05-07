# Task 7: Swagger Documentation

## Alur Kerja

1. Sistem menggunakan Swagger untuk mendokumentasikan API
2. Dokumentasi API dibuat menggunakan komentar anotasi di kode sumber
3. Spesifikasi API dihasilkan secara otomatis menggunakan perintah `swag init`
4. Dokumentasi API mencakup semua endpoint, parameter, request body, dan response
5. Swagger UI diintegrasikan ke dalam aplikasi untuk memudahkan pengujian API
6. Dokumentasi API diperbarui setiap kali ada perubahan pada endpoint

## Endpoint API

- **GET /swagger/index.html**
  - **Deskripsi**: Menampilkan dokumentasi Swagger UI interaktif
  - **Response**: Halaman HTML Swagger UI

- **GET /swagger/doc.json**
  - **Deskripsi**: Mengembalikan spesifikasi OpenAPI dalam format JSON
  - **Response**: File JSON yang berisi spesifikasi OpenAPI

## Pengujian

1. Jalankan aplikasi:
   ```bash
   docker-compose up -d
   ```

2. Akses Swagger UI melalui browser:
   ```
   http://localhost:8080/swagger/index.html
   ```