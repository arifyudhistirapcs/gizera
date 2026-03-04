#!/bin/bash

echo "Testing Attendance API Endpoints"
echo "================================="
echo ""

# Login first
echo "1. Login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier": "TEST001", "password": "password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "✗ Login failed"
  exit 1
fi

echo "✓ Login successful"
echo ""

# Test today's attendance
echo "2. Get today's attendance..."
curl -s -X GET http://localhost:8080/api/v1/attendance/today \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo ""

# Test attendance history (7 days)
echo "3. Get attendance history (7 days)..."
END_DATE=$(date +%Y-%m-%d)
START_DATE=$(date -v-7d +%Y-%m-%d 2>/dev/null || date -d '7 days ago' +%Y-%m-%d)

echo "Date range: $START_DATE to $END_DATE"
curl -s -X GET "http://localhost:8080/api/v1/attendance/by-date-range?start_date=$START_DATE&end_date=$END_DATE" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo ""

echo "================================="
echo "Test complete!"
