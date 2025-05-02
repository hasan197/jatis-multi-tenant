#!/bin/sh

# Menjalankan migrasi database di dalam container
echo "Running database migrations in container..."

# Menunggu PostgreSQL siap
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing migrations"

# Koneksi ke PostgreSQL dan menjalankan skrip SQL
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f /app/scripts/init_db.sql

# Menjalankan semua file migration secara berurutan
for migration in /app/scripts/migrations/*.sql; do
    echo "Running migration: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$migration"
done

echo "Migrations completed!" 