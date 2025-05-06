-- Add workers column to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS workers INTEGER NOT NULL DEFAULT 3;
