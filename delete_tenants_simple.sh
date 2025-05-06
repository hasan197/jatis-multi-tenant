#!/bin/bash

# Ambil semua tenant
TENANTS=$(curl -s http://localhost:8080/api/tenants)
echo "Daftar tenant yang akan dihapus:"
echo "$TENANTS"
echo "-------------------------------------"

# Ekstrak ID tenant menggunakan grep dan cut (lebih sederhana dari jq)
echo "$TENANTS" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | while read id; do
  echo "Menghapus tenant dengan ID: $id"
  curl -X DELETE http://localhost:8080/api/tenants/$id
  echo -e "\n"
done

echo "Semua tenant telah dihapus."
