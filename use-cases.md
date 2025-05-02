# Use Cases Sistem Messaging Multi-Tenant

## 1. Manajemen Tenant

### 1.1 Pembuatan Tenant Baru
**Aktor**: Administrator Sistem
**Deskripsi**: Administrator membuat tenant baru dalam sistem
**Alur**:
1. Administrator mengirim request POST ke `/tenants`
2. Sistem membuat queue RabbitMQ baru dengan format `tenant_{id}_queue`
3. Sistem memulai consumer goroutine untuk tenant tersebut
4. Sistem menyimpan informasi tenant dan consumer di TenantManager
5. Sistem mengembalikan response sukses dengan ID tenant

### 1.2 Penghapusan Tenant
**Aktor**: Administrator Sistem
**Deskripsi**: Administrator menghapus tenant yang ada
**Alur**:
1. Administrator mengirim request DELETE ke `/tenants/{id}`
2. Sistem mengirim sinyal shutdown ke consumer tenant
3. Sistem menutup channel RabbitMQ dan menghapus queue
4. Sistem menghapus tenant dari TenantManager
5. Sistem mengembalikan response sukses

## 2. Pengiriman dan Penerimaan Pesan

### 2.1 Pengiriman Pesan
**Aktor**: Aplikasi Tenant
**Deskripsi**: Tenant mengirim pesan ke sistem
**Alur**:
1. Aplikasi tenant mengirim pesan ke queue tenant
2. Consumer tenant menerima pesan
3. Sistem menyimpan pesan ke tabel messages yang terpartisi
4. Sistem mengirim acknowledgment ke RabbitMQ

### 2.2 Penerimaan Pesan
**Aktor**: Aplikasi Tenant
**Deskripsi**: Tenant menerima pesan dari sistem
**Alur**:
1. Aplikasi tenant melakukan polling ke endpoint GET `/messages`
2. Sistem mengambil pesan berdasarkan cursor pagination
3. Sistem mengembalikan pesan dan cursor berikutnya

## 3. Konfigurasi Worker

### 3.1 Update Konfigurasi Worker
**Aktor**: Administrator Tenant
**Deskripsi**: Administrator mengubah jumlah worker untuk tenant
**Alur**:
1. Administrator mengirim request PUT ke `/tenants/{id}/config/concurrency`
2. Sistem memvalidasi jumlah worker yang diminta
3. Sistem mengupdate jumlah worker secara atomik
4. Sistem menyesuaikan worker pool sesuai konfigurasi baru
5. Sistem mengembalikan response sukses

## 4. Graceful Shutdown

### 4.1 Shutdown Aplikasi
**Aktor**: Administrator Sistem
**Deskripsi**: Administrator melakukan shutdown aplikasi
**Alur**:
1. Administrator mengirim sinyal shutdown ke aplikasi
2. Sistem menghentikan penerimaan request baru
3. Sistem menunggu semua transaksi yang sedang berjalan selesai
4. Sistem menutup koneksi ke RabbitMQ dan PostgreSQL
5. Sistem menghentikan aplikasi dengan bersih

## 5. Monitoring dan Maintenance

### 5.1 Monitoring Queue
**Aktor**: Administrator Sistem
**Deskripsi**: Administrator memonitor status queue
**Alur**:
1. Sistem mengumpulkan metrik Prometheus untuk queue depth
2. Sistem menampilkan metrik melalui dashboard monitoring
3. Administrator dapat melihat status queue setiap tenant

### 5.2 Penanganan Pesan Gagal
**Aktor**: Sistem
**Deskripsi**: Sistem menangani pesan yang gagal diproses
**Alur**:
1. Sistem mendeteksi pesan yang gagal diproses
2. Sistem memindahkan pesan ke dead-letter queue
3. Sistem mencoba memproses ulang pesan sesuai kebijakan retry
4. Sistem mencatat log untuk pesan yang gagal

## 6. Keamanan

### 6.1 Autentikasi Tenant
**Aktor**: Aplikasi Tenant
**Deskripsi**: Tenant melakukan autentikasi untuk mengakses sistem
**Alur**:
1. Tenant mengirim kredensial untuk autentikasi
2. Sistem memvalidasi JWT token
3. Sistem memberikan akses ke resource tenant yang sesuai
4. Sistem menolak akses ke resource tenant lain 