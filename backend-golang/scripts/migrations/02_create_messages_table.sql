-- Membuat tabel messages dengan partisi berdasarkan tenant_id
CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, id)
) PARTITION BY LIST (tenant_id);

-- Membuat index pada tenant_id dan created_at
CREATE INDEX IF NOT EXISTS idx_messages_tenant_id ON messages(tenant_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Membuat fungsi untuk membuat partisi baru
CREATE OR REPLACE FUNCTION create_messages_partition(tenant_id UUID)
RETURNS void AS $$
DECLARE
    partition_name TEXT;
BEGIN
    partition_name := 'messages_' || replace(tenant_id::text, '-', '_');
    
    -- Membuat partisi baru jika belum ada
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF messages FOR VALUES IN (%L)',
        partition_name,
        tenant_id
    );
END;
$$ LANGUAGE plpgsql;

-- Membuat trigger untuk memastikan partisi ada sebelum insert
CREATE OR REPLACE FUNCTION ensure_messages_partition()
RETURNS trigger AS $$
BEGIN
    PERFORM create_messages_partition(NEW.tenant_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_ensure_messages_partition
    BEFORE INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION ensure_messages_partition(); 