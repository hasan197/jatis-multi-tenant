#!/bin/bash

# Menjalankan migrasi database dari host
echo "Running database migrations..."

# Menjalankan migration di dalam container menggunakan nerdctl
nerdctl compose exec backend-golang sh -c "cd /app && ./scripts/run_migrations_in_container.sh"

echo "Migrations completed!" 