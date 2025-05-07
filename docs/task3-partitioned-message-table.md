# Task 3: Partitioned Message Table

## Alur Kerja

1. Sistem menggunakan partisi database untuk memisahkan penyimpanan pesan berdasarkan tenant
2. Ketika tenant baru dibuat:
   - Sistem membuat partisi tabel baru dengan format `messages_{tenant_id}`
   - Partisi ini terhubung ke tabel utama `messages` melalui partisi LIST berdasarkan `tenant_id`
3. Ketika pesan di insert ke tabel `messages`:
   - PostgreSQL secara otomatis menyimpan data ke partisi yang sesuai berdasarkan `tenant_id`
4. Ketika tenant dihapus:
   - Sistem menghapus partisi tabel yang terkait dengan tenant tersebut
   - Data pesan tenant tersebut juga ikut terhapus

## Endpoint API

Tidak ada endpoint API khusus untuk fitur ini. Partisi tabel dibuat dan dikelola secara otomatis saat tenant dibuat atau dihapus melalui endpoint tenant.

Skema database yang digunakan:
```sql
CREATE TABLE messages (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  payload JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY LIST (tenant_id);
```

## Pengujian

1. Buat tenant baru:
   ```bash
   curl -X POST http://localhost:8080/api/tenants -H "Content-Type: application/json" -d '{"name":"Test Partition","description":"Test tenant for partition verification","status":"active"}'
   ```

2. Verifikasi partisi tabel telah dibuat:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT tablename FROM pg_tables WHERE tablename LIKE 'messages_%';"
   ```

3. Kirim beberapa pesan ke tenant tersebut dan verifikasi pesan disimpan di partisi yang benar:
   ```bash
   nerdctl exec -it jatis-sample-stack-golang-postgres-1 psql -U postgres -d sample_db -c "SELECT count(*) FROM messages_{tenant_id};"
   ```
