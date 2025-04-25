#!/bin/bash

# Menjalankan migrasi database
echo "Running database migrations..."

# Koneksi ke PostgreSQL dan menjalankan skrip SQL
psql -h postgres -U postgres -d sample_db -f /app/scripts/init_db.sql

echo "Migrations completed!" 