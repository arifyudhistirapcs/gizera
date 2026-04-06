#!/bin/bash

# Test script for IP-based check-in
# This script tests the check-in functionality with IP validation

echo "=== Testing IP-Based Check-in ==="
echo ""

# Step 1: Login to get token
echo "Step 1: Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST001",
    "password": "password123"
  }')

echo "Login Response: $LOGIN_RESPONSE"
echo ""

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "ERROR: Failed to get token"
  exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."
echo ""

# Step 2: Check current attendance
echo "Step 2: Checking today's attendance..."
TODAY_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/attendance/today \
  -H "Authorization: Bearer $TOKEN")

echo "Today's Attendance: $TODAY_RESPONSE"
echo ""

# Step 3: Perform check-in with IP validation
echo "Step 3: Performing check-in (IP-based validation)..."
CHECKIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "ssid": "AUTO-DETECT",
    "bssid": "00:00:00:00:00:00"
  }')

echo "Check-in Response: $CHECKIN_RESPONSE"
echo ""

# Check if successful
if echo "$CHECKIN_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Check-in SUCCESSFUL!"
  
  # Extract validated_by info
  echo ""
  echo "Validation Details:"
  echo "$CHECKIN_RESPONSE" | grep -o '"validated_by":{[^}]*}' || echo "No validation details"
else
  echo "❌ Check-in FAILED!"
  
  # Show error details
  ERROR_CODE=$(echo "$CHECKIN_RESPONSE" | grep -o '"error_code":"[^"]*' | cut -d'"' -f4)
  MESSAGE=$(echo "$CHECKIN_RESPONSE" | grep -o '"message":"[^"]*' | cut -d'"' -f4)
  CLIENT_IP=$(echo "$CHECKIN_RESPONSE" | grep -o '"client_ip":"[^"]*' | cut -d'"' -f4)
  
  echo "Error Code: $ERROR_CODE"
  echo "Message: $MESSAGE"
  echo "Client IP: $CLIENT_IP"
fi

echo ""
echo "=== Test Complete ==="
