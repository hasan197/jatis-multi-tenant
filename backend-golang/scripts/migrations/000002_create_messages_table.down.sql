-- Drop trigger
DROP TRIGGER IF EXISTS trigger_ensure_messages_partition ON messages;

-- Drop functions
DROP FUNCTION IF EXISTS ensure_messages_partition();
DROP FUNCTION IF EXISTS create_messages_partition(UUID);

-- Drop messages table
DROP TABLE IF EXISTS messages; 