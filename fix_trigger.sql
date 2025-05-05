DROP TRIGGER IF EXISTS trigger_ensure_messages_partition ON messages;
DROP FUNCTION IF EXISTS ensure_messages_partition();

CREATE OR REPLACE FUNCTION ensure_messages_partition()
RETURNS trigger AS $$
BEGIN
    EXECUTE format('CREATE TABLE IF NOT EXISTS messages_%s PARTITION OF messages FOR VALUES IN (%L)',
        replace(NEW.tenant_id::text, '-', '_'),
        NEW.tenant_id
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_ensure_messages_partition
    BEFORE INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION ensure_messages_partition(); 