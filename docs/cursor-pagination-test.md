# Panduan Pengujian Cursor Pagination API

Dokumen ini menjelaskan cara menguji fitur cursor pagination API pada aplikasi backend-golang.

## Prasyarat

- curl untuk mengirim request HTTP
- Terminal/command prompt
- Aplikasi backend-golang yang berjalan di http://localhost:8080

## Langkah-langkah Pengujian

### 1. Memastikan Aplikasi Berjalan

```bash
# Periksa status container
nerdctl ps | grep backend-golang

# Jika container tidak berjalan, jalankan dengan docker-compose
nerdctl compose up -d
```

### 2. Membuat Tenant Baru untuk Pengujian

```bash
# Buat tenant baru untuk pengujian
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"Pagination Test","description":"Tenant untuk menguji fitur cursor pagination","status":"active","workers":3}' \
  http://localhost:8080/api/tenants

# Catat ID tenant yang dikembalikan untuk digunakan pada langkah selanjutnya
# Contoh output: {"id":"b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",...}
```

### 3. Menambahkan Beberapa Pesan untuk Pengujian

```bash
# Tambahkan pesan pertama
curl -X POST -H "Content-Type: application/json" \
  -d '{"payload": {"message": "Pesan pertama untuk pengujian cursor pagination", "priority": "high"}}' \
  http://localhost:8080/api/tenants/TENANT_ID/messages

# Tambahkan pesan kedua
curl -X POST -H "Content-Type: application/json" \
  -d '{"payload": {"message": "Pesan kedua untuk pengujian cursor pagination", "priority": "medium"}}' \
  http://localhost:8080/api/tenants/TENANT_ID/messages

# Tambahkan pesan ketiga
curl -X POST -H "Content-Type: application/json" \
  -d '{"payload": {"message": "Pesan ketiga untuk pengujian cursor pagination", "priority": "low"}}' \
  http://localhost:8080/api/tenants/TENANT_ID/messages

# Tambahkan pesan keempat
curl -X POST -H "Content-Type: application/json" \
  -d '{"payload": {"message": "Pesan keempat untuk pengujian cursor pagination", "priority": "high"}}' \
  http://localhost:8080/api/tenants/TENANT_ID/messages

# Tambahkan pesan kelima
curl -X POST -H "Content-Type: application/json" \
  -d '{"payload": {"message": "Pesan kelima untuk pengujian cursor pagination", "priority": "medium"}}' \
  http://localhost:8080/api/tenants/TENANT_ID/messages
```

### 4. Menggunakan Endpoint Cursor Pagination

#### 4.1. Mengambil Halaman Pertama (Tanpa Cursor)

```bash
# Mengambil 2 pesan pertama
curl -s "http://localhost:8080/api/messages?limit=2"

# Contoh output:
# {
#   "data": [
#     {
#       "id": "43df1ac1-c50f-4825-993c-d178f2b84f87",
#       "tenant_id": "b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",
#       "payload": {
#         "message": "Pesan pertama untuk pengujian cursor pagination",
#         "priority": "high"
#       },
#       "created_at": "2025-05-06T11:16:52.759268Z",
#       "updated_at": "2025-05-06T11:16:52.759294Z"
#     },
#     {
#       "id": "b1ad3e1a-b7b4-417b-b878-4531f8e68cfb",
#       "tenant_id": "b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",
#       "payload": {
#         "message": "Pesan ketiga untuk pengujian cursor pagination",
#         "priority": "low"
#       },
#       "created_at": "2025-05-06T11:17:19.149183Z",
#       "updated_at": "2025-05-06T11:17:19.14923Z"
#     }
#   ],
#   "next_cursor": "b1ad3e1a-b7b4-417b-b878-4531f8e68cfb"
# }

# Catat nilai next_cursor untuk digunakan pada langkah selanjutnya
```

#### 4.2. Mengambil Halaman Berikutnya (Dengan Cursor)

```bash
# Mengambil halaman berikutnya menggunakan cursor dari langkah sebelumnya
curl -s "http://localhost:8080/api/messages?cursor=NEXT_CURSOR&limit=2"

# Contoh:
curl -s "http://localhost:8080/api/messages?cursor=b1ad3e1a-b7b4-417b-b878-4531f8e68cfb&limit=2"

# Contoh output:
# {
#   "data": [
#     {
#       "id": "c24c4ac4-fe75-404a-9e9f-6ba73caa2825",
#       "tenant_id": "b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",
#       "payload": {
#         "message": "Pesan kedua untuk pengujian cursor pagination",
#         "priority": "medium"
#       },
#       "created_at": "2025-05-06T11:17:06.399618Z",
#       "updated_at": "2025-05-06T11:17:06.399642Z"
#     },
#     {
#       "id": "e5df2bc2-d61f-4936-aa4c-e289f3c95f98",
#       "tenant_id": "b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",
#       "payload": {
#         "message": "Pesan keempat untuk pengujian cursor pagination",
#         "priority": "high"
#       },
#       "created_at": "2025-05-06T11:18:30.149183Z",
#       "updated_at": "2025-05-06T11:18:30.14923Z"
#     }
#   ],
#   "next_cursor": "e5df2bc2-d61f-4936-aa4c-e289f3c95f98"
# }
```

#### 4.3. Mengambil Halaman Terakhir

```bash
# Mengambil halaman terakhir menggunakan cursor dari langkah sebelumnya
curl -s "http://localhost:8080/api/messages?cursor=NEXT_CURSOR&limit=2"

# Contoh:
curl -s "http://localhost:8080/api/messages?cursor=e5df2bc2-d61f-4936-aa4c-e289f3c95f98&limit=2"

# Contoh output (halaman terakhir dengan next_cursor kosong):
# {
#   "data": [
#     {
#       "id": "f6eg3cd3-e72g-5047-bb5d-f390g4d06g09",
#       "tenant_id": "b1adc5c3-0878-44f1-9ad8-5c1ad40c6879",
#       "payload": {
#         "message": "Pesan kelima untuk pengujian cursor pagination",
#         "priority": "medium"
#       },
#       "created_at": "2025-05-06T11:19:45.149183Z",
#       "updated_at": "2025-05-06T11:19:45.14923Z"
#     }
#   ],
#   "next_cursor": ""
# }
```

### 5. Mengubah Limit untuk Mengambil Lebih Banyak atau Lebih Sedikit Data

```bash
# Mengambil 1 pesan per halaman
curl -s "http://localhost:8080/api/messages?limit=1"

# Mengambil 5 pesan per halaman
curl -s "http://localhost:8080/api/messages?limit=5"

# Mengambil semua pesan (maksimum 100)
curl -s "http://localhost:8080/api/messages?limit=100"
```

## Verifikasi Hasil

Fitur cursor pagination berfungsi dengan benar jika:

1. API mengembalikan response dalam format yang sesuai dengan spesifikasi:
   ```json
   {
     "data": [...],
     "next_cursor": "456"
   }
   ```

2. Jumlah item dalam array `data` sesuai dengan parameter `limit` yang diberikan (kecuali pada halaman terakhir yang mungkin memiliki jumlah item lebih sedikit)

3. Nilai `next_cursor` berisi ID dari item terakhir di halaman saat ini, yang dapat digunakan untuk mengambil halaman berikutnya

4. Nilai `next_cursor` kosong (`""`) pada halaman terakhir, menandakan bahwa tidak ada lagi data yang tersedia

## Catatan Penting

- Pastikan untuk mengganti `TENANT_ID` dengan ID tenant yang sebenarnya dari hasil langkah 2
- Pastikan untuk mengganti `NEXT_CURSOR` dengan nilai next_cursor yang diperoleh dari response sebelumnya
- Parameter `limit` dapat diatur antara 1 hingga 100, dengan nilai default 10 jika tidak ditentukan
- Cursor pagination menjamin tidak ada data yang hilang atau duplikat saat navigasi antar halaman, bahkan jika data baru ditambahkan selama proses pagination
