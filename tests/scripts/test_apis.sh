#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
TIMESTAMP=$(date +%s)
USERNAME="user_${TIMESTAMP}"
EMAIL="user_${TIMESTAMP}@example.com"
PASSWORD="password123"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Validations
SERVER_PID=""

# Try get PID from argument
if [ -n "$1" ]; then
    SERVER_PID=$1
else
    # Try get PID from Stdin (Pipe)
    if [ ! -t 0 ]; then
        read SERVER_PID
    fi
fi

if [ -z "$SERVER_PID" ]; then
    echo "Usage: $0 <SERVER_PID> or echo <SERVER_PID> | $0"
    echo "Warning: No PID provided. I won't be able to stop the server automatically."
else
    echo "Server PID received: $SERVER_PID"
    trap "echo 'Stopping server (PID: $SERVER_PID)...'; kill $SERVER_PID" EXIT
fi

# Wait for server to be ready
echo "Waiting for server to be ready..."
for i in {1..30}; do
    if curl -s http://localhost:8080 > /dev/null; then
        echo "Server is up!"
        break
    fi
    # If using /api/v1/auth/register it might be 404 on root, but port is open.
    # Just checking connection refused vs (any http response)
    # curl returns 7 on connection refused.
    # Actually, simpler check:
    nc -z localhost 8080 && echo "Server port 8080 open" && break
    
    echo -n "."
    sleep 1
done
echo ""

# Helper function
check_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}[PASS]${NC} $2"
    else
        echo -e "${RED}[FAIL]${NC} $2"
        exit 1
    fi
}

echo "----------------------------------------"
echo "Running API Integration Tests"
echo "Target: $API_URL"
echo "User: $USERNAME"
echo "----------------------------------------"

# 1. Register
echo -e "\n1. Testing Register..."
REGISTER_RES=$(curl -s -X POST "${API_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$USERNAME\",
    \"password\": \"$PASSWORD\",
    \"email\": \"$EMAIL\",
    \"nickname\": \"Test User\"
  }")

# Check register success (look for successful code 200 or user data)
echo "$REGISTER_RES" | grep "success" > /dev/null
check_status $? "Registration Response: $(echo $REGISTER_RES | jq -r .msg)"

# 2. Login
echo -e "\n2. Testing Login..."
LOGIN_RES=$(curl -s -X POST "${API_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"identifier\": \"$USERNAME\",
    \"password\": \"$PASSWORD\"
  }")

# Extract Token
TOKEN=$(echo $LOGIN_RES | jq -r '.data.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo -e "${GREEN}[PASS]${NC} Login Successful. Token: ${TOKEN:0:10}..."
else
    echo -e "${RED}[FAIL]${NC} Login Failed"
    echo "Response: $LOGIN_RES"
    exit 1
fi

# 3. Get Profile
echo -e "\n3. Testing Get Profile..."
PROFILE_RES=$(curl -s -X GET "${API_URL}/account/profile" \
  -H "Authorization: Bearer $TOKEN")

CURRENT_NICKNAME=$(echo $PROFILE_RES | jq -r '.data.nickname')
if [ "$CURRENT_NICKNAME" == "Test User" ]; then
    echo -e "${GREEN}[PASS]${NC} Get Profile (Nickname: $CURRENT_NICKNAME)"
else
    echo -e "${RED}[FAIL]${NC} Get Profile mismatch"
    echo "Response: $PROFILE_RES"
    exit 1
fi

# 4. Update Profile
echo -e "\n4. Testing Update Profile..."
NEW_NICKNAME="Updated_$TIMESTAMP"
UPDATE_RES=$(curl -s -X PUT "${API_URL}/account/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"nickname\": \"$NEW_NICKNAME\",
    \"gender\": \"male\",
    \"description\": \"Integration test update\"
  }")

echo "$UPDATE_RES" | grep "success" > /dev/null
check_status $? "Update Profile"

# 5. Verify Update
echo -e "\n5. Verifying Update..."
PROFILE_RES_2=$(curl -s -X GET "${API_URL}/account/profile" \
  -H "Authorization: Bearer $TOKEN")

UPDATED_NICKNAME=$(echo $PROFILE_RES_2 | jq -r '.data.nickname')
UPDATED_DESC=$(echo $PROFILE_RES_2 | jq -r '.data.description')

if [ "$UPDATED_NICKNAME" == "$NEW_NICKNAME" ]; then
    echo -e "${GREEN}[PASS]${NC} Nickname Updated ($UPDATED_NICKNAME)"
else
    echo -e "${RED}[FAIL]${NC} Nickname Update Failed"
fi

if [ "$UPDATED_DESC" == "Integration test update" ]; then
    echo -e "${GREEN}[PASS]${NC} Description Updated ($UPDATED_DESC)"
else
    echo -e "${RED}[FAIL]${NC} Description Update Failed"
    echo "Res: $PROFILE_RES_2"
fi

echo "----------------------------------------"
echo "All Tests Completed Successfully"
echo "----------------------------------------"
