-- Hapus fungsi-fungsi
DROP FUNCTION IF EXISTS maintain_messages_partitions();
DROP FUNCTION IF EXISTS check_messages_partition(UUID);
DROP FUNCTION IF EXISTS drop_messages_partition(UUID);
DROP FUNCTION IF EXISTS ensure_messages_partition();
DROP FUNCTION IF EXISTS create_messages_partition(UUID);

-- Hapus trigger
DROP TRIGGER IF EXISTS trigger_ensure_messages_partition ON messages;

-- Hapus index
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_tenant_id;

-- Hapus tabel messages dan semua partisinya
DROP TABLE IF EXISTS messages CASCADE; 