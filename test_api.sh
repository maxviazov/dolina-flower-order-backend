#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Testing Dolina Flower Order Backend API ==="
echo ""

echo "1. Health Check"
curl -s "$BASE_URL/health" | jq .
echo -e "\n"

echo "2. Ping"
curl -s "$BASE_URL/api/v1/ping" | jq .
echo -e "\n"

echo "3. Get Available Flowers"
curl -s "$BASE_URL/api/v1/flowers" | jq .
echo -e "\n"

echo "4. Create Order"
curl -s -X POST "$BASE_URL/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "mark_box": "VVA",
    "customer_id": "test-customer-123",
    "items": [
      {
        "variety": "Red Naomi",
        "length": 70,
        "box_count": 10.5,
        "pack_rate": 20,
        "total_stems": 210,
        "farm_name": "KENYA FARM 1",
        "truck_name": "TRUCK A",
        "price": 4.07
      }
    ],
    "notes": "Test order"
  }' | jq .
echo -e "\n"

echo "Done!"
