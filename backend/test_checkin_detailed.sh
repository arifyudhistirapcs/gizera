#!/bin/bash

echo "========================================="
echo "  IP-Based Check-in Validation Test"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Login
echo -e "${YELLOW}[1/4] Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "TEST001",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo -e "${RED}✗ Login failed${NC}"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo ""

# Step 2: Check WiFi Config
echo -e "${YELLOW}[2/4] Checking WiFi configuration...${NC}"
WIFI_CONFIG=$(curl -s -X GET "http://localhost:8080/api/v1/wifi-config?active_only=true" \
  -H "Authorization: Bearer $TOKEN")

echo "WiFi Configs:"
echo "$WIFI_CONFIG" | python3 -m json.tool 2>/dev/null || echo "$WIFI_CONFIG"
echo ""

# Step 3: Check current attendance
echo -e "${YELLOW}[3/4] Checking today's attendance...${NC}"
TODAY_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/attendance/today \
  -H "Authorization: Bearer $TOKEN")

echo "Today's Attendance:"
echo "$TODAY_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$TODAY_RESPONSE"
echo ""

# Step 4: Perform check-in
echo -e "${YELLOW}[4/4] Performing check-in...${NC}"
echo "Request payload: {ssid: 'AUTO-DETECT', bssid: '00:00:00:00:00:00'}"
echo ""

CHECKIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "ssid": "AUTO-DETECT",
    "bssid": "00:00:00:00:00:00"
  }')

echo "Check-in Response:"
echo "$CHECKIN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CHECKIN_RESPONSE"
echo ""

# Parse response
if echo "$CHECKIN_RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}========================================="
  echo -e "  ✓ CHECK-IN SUCCESSFUL"
  echo -e "=========================================${NC}"
  
  # Extract validation info
  METHOD=$(echo "$CHECKIN_RESPONSE" | grep -o '"method":"[^"]*' | cut -d'"' -f4)
  CLIENT_IP=$(echo "$CHECKIN_RESPONSE" | grep -o '"client_ip":"[^"]*' | cut -d'"' -f4)
  
  echo ""
  echo "Validation Details:"
  echo "  Method: $METHOD"
  echo "  Client IP: $CLIENT_IP"
  echo ""
  
  if [ "$METHOD" = "ip_validation" ]; then
    echo -e "${GREEN}✓ Validated via IP address${NC}"
  else
    echo -e "${YELLOW}⚠ Validated via SSID/BSSID (fallback)${NC}"
  fi
  
else
  echo -e "${RED}========================================="
  echo -e "  ✗ CHECK-IN FAILED"
  echo -e "=========================================${NC}"
  
  ERROR_CODE=$(echo "$CHECKIN_RESPONSE" | grep -o '"error_code":"[^"]*' | cut -d'"' -f4)
  MESSAGE=$(echo "$CHECKIN_RESPONSE" | grep -o '"message":"[^"]*' | cut -d'"' -f4)
  CLIENT_IP=$(echo "$CHECKIN_RESPONSE" | grep -o '"client_ip":"[^"]*' | cut -d'"' -f4)
  
  echo ""
  echo "Error Details:"
  echo "  Error Code: $ERROR_CODE"
  echo "  Message: $MESSAGE"
  echo "  Client IP: $CLIENT_IP"
fi

echo ""
echo "========================================="
echo "  Check backend logs for detailed info"
echo "========================================="
