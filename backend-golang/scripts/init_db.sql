-- Membuat tabel users jika belum ada
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Membuat index pada email untuk mempercepat pencarian
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Menambahkan data awal (password = "password123" yang sudah di-hash)
INSERT INTO users (name, email, password, created_at, updated_at)
VALUES 
    ('Admin User', 'admin@example.com', '$2a$10$zDVKMIYiALxUoIpHgZ9.l.lZPL/wn.U4OxLW9X3Y0z0aeYChR.s1G', NOW(), NOW()),
    ('Test User', 'user@example.com', '$2a$10$zDVKMIYiALxUoIpHgZ9.l.lZPL/wn.U4OxLW9X3Y0z0aeYChR.s1G', NOW(), NOW())
ON CONFLICT (email) DO NOTHING; 