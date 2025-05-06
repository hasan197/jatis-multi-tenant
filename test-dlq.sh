#!/bin/bash

# Script untuk menguji Dead Letter Queue

echo "=== Memulai pengujian Dead Letter Queue ==="

# 1. Membuat tenant baru dengan ID unik
TENANT_ID="test-dlq-$(date +%s)"
echo "Menggunakan tenant ID: $TENANT_ID"

echo "1. Membuat tenant baru dengan ID: $TENANT_ID"
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/tenants" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Tenant DLQ","workers":2}')

echo "$RESPONSE"

# Ekstrak tenant ID dari respons jika ada
if [[ "$RESPONSE" == *"id"* ]]; then
  # Ekstrak ID dari respons JSON (tanpa jq untuk kompatibilitas)
  TENANT_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | cut -d '"' -f 4)
  echo "Menggunakan tenant ID dari respons: $TENANT_ID"
fi

# 2. Menunggu sebentar agar tenant terdaftar dan consumer aktif
echo "2. Menunggu tenant terdaftar dan consumer aktif..."
sleep 3

# 3. Mempublikasikan pesan yang akan gagal diproses (dengan error)
echo "3. Mempublikasikan pesan yang akan gagal (dengan error 'force_error')"
PUBLISH_RESPONSE=$(curl -s -X POST "http://localhost:3000/api/tenants/$TENANT_ID/publish" \
  -H "Content-Type: application/json" \
  -d '{"content":"Test message with error","metadata":{"force_error":true}}')

echo "$PUBLISH_RESPONSE"

# 4. Menunggu beberapa saat agar pesan diproses dan retry dilakukan
echo "4. Menunggu pesan diproses dan retry dilakukan (15 detik)..."
sleep 5
echo "   Masih menunggu retry (10 detik lagi)..."
sleep 5
echo "   Menunggu retry terakhir (5 detik lagi)..."
sleep 5

# 5. Memeriksa status queue untuk melihat apakah pesan masuk ke DLQ
echo "5. Memeriksa status queue untuk tenant $TENANT_ID"
QUEUE_STATUS=$(curl -s -X GET "http://localhost:8080/api/tenants/$TENANT_ID/queue-status")
echo "$QUEUE_STATUS"

# 6. Memeriksa pesan di database untuk melihat status pemrosesan
echo "6. Memeriksa pesan di database"
MESSAGES=$(curl -s -X GET "http://localhost:8080/api/messages?tenant_id=$TENANT_ID&limit=10")
echo "$MESSAGES"

# 7. Memeriksa log dari backend-golang
echo "\n7. Memeriksa log dari backend-golang untuk tenant $TENANT_ID dan DLQ"
nerdctl logs jatis-sample-stack-golang-backend-golang-1 | grep -E "$TENANT_ID|retry|dead-letter|DLQ|error" | tail -n 30

echo "\n=== Pengujian Dead Letter Queue selesai ==="
