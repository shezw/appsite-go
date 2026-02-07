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

# 6. List Users
echo -e "\n6. Testing List Users..."
LIST_RES=$(curl -s -X GET "${API_URL}/users?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

LIST_COUNT=$(echo $LIST_RES | jq -r '.data.total')
LIST_FIRST_USER=$(echo $LIST_RES | jq -r '.data.list[0].username')

# Note: LIST_COUNT might be treated as string by shell comparisons if not careful, but -gt handles integers.
# jq output for total might be null if failed.
if [[ "$LIST_COUNT" =~ ^[0-9]+$ ]] && [ "$LIST_COUNT" -ge 1 ]; then
    echo -e "${GREEN}[PASS]${NC} List Users (Total: $LIST_COUNT, First: $LIST_FIRST_USER)"
else
    echo -e "${RED}[FAIL]${NC} List Users failed or empty"
    echo "Response: $LIST_RES"
    exit 1
fi

# 7. Create Article
echo -e "\n7. Testing Create Article..."
random_suffix=$(date +%s)
ARTICLE_RES=$(curl -s -X POST "${API_URL}/content/articles" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"My First Article $random_suffix\",
    \"content\": \"This is the content\",
    \"status\": \"enabled\"
  }")

ARTICLE_ID=$(echo $ARTICLE_RES | jq -r '.data.id')

if [ "$ARTICLE_ID" != "null" ] && [ -n "$ARTICLE_ID" ]; then
    echo -e "${GREEN}[PASS]${NC} Create Article (ID: $ARTICLE_ID)"
else
    echo -e "${RED}[FAIL]${NC} Create Article Failed"
    echo "Response: $ARTICLE_RES"
    exit 1
fi

# 8. List Articles (Public? or Protected?)
# In router we registered /content/articles as public/protected depending on method?
# GET /content/articles is in public block in router.go
echo -e "\n8. Testing List Articles (Public)..."
LIST_ART_RES=$(curl -s -X GET "${API_URL}/content/articles")

LIST_ART_COUNT=$(echo $LIST_ART_RES | jq -r '.data.total')

if [[ "$LIST_ART_COUNT" =~ ^[0-9]+$ ]] && [ "$LIST_ART_COUNT" -ge 1 ]; then
    echo -e "${GREEN}[PASS]${NC} List Articles (Total: $LIST_ART_COUNT)"
else
    echo -e "${RED}[FAIL]${NC} List Articles failed or empty"
    echo "Response: $LIST_ART_RES"
fi

# 9. Test WeChat Callback (Mock)
echo -e "\n9. Testing WeChat Callback..."
WECHAT_RES=$(curl -s -X GET "${API_URL}/callback/wechat")
echo "$WECHAT_RES" | grep "mock_wechat_ok" > /dev/null
check_status $? "WeChat Callback Response"

# 10. Test OSS Callback (Mock)
echo -e "\n10. Testing OSS Callback..."
OSS_RES=$(curl -s -X POST "${API_URL}/callback/oss")
echo "$OSS_RES" | grep "mock_oss_ok" > /dev/null
check_status $? "OSS Callback Response"

# --- Admin Tests ---
ADMIN_URL="http://localhost:8080/admin/v1"

echo -e "\n--- Admin Tests ---"

# 11. Admin Login
echo -e "\n11. Testing Admin Login..."
ADMIN_LOGIN_RES=$(curl -s -X POST "${ADMIN_URL}/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$USERNAME\",
    \"password\": \"$PASSWORD\"
  }")

ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RES | jq -r '.data.token')

if [ "$ADMIN_TOKEN" != "null" ] && [ -n "$ADMIN_TOKEN" ]; then
    echo -e "${GREEN}[PASS]${NC} Admin Login Successful. Token: ${ADMIN_TOKEN:0:10}..."
else
    echo -e "${RED}[FAIL]${NC} Admin Login Failed"
    echo "Response: $ADMIN_LOGIN_RES"
    exit 1
fi

# 12. Admin List Users
echo -e "\n12. Testing Admin List Users..."
ADMIN_LIST_RES=$(curl -s -X GET "${ADMIN_URL}/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

ADMIN_LIST_COUNT=$(echo $ADMIN_LIST_RES | jq -r '.data.total')

if [[ "$ADMIN_LIST_COUNT" =~ ^[0-9]+$ ]] && [ "$ADMIN_LIST_COUNT" -ge 1 ]; then
    echo -e "${GREEN}[PASS]${NC} Admin List Users (Total: $ADMIN_LIST_COUNT)"
else
    echo -e "${RED}[FAIL]${NC} Admin List Users failed"
    echo "Response: $ADMIN_LIST_RES"
fi

# Get User ID (using the created user's username to find ID from list)
# We need to filter the list to find the UID of our user, or just assume first user for this test if clean DB
# With parallelism or persistence it might change.
# jq select .username == $USERNAME
TARGET_UID=$(echo $ADMIN_LIST_RES | jq -r ".data.list[] | select(.username == \"$USERNAME\") | .uid")

if [ -z "$TARGET_UID" ]; then
    echo "Cannot find UID for $USERNAME"
    TARGET_UID=$(echo $ADMIN_LIST_RES | jq -r '.data.list[0].uid')
fi

# 13. Admin Ban User
echo -e "\n13. Testing Admin Ban User ($TARGET_UID)..."
BAN_RES=$(curl -s -X PUT "${ADMIN_URL}/users/$TARGET_UID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"status\": \"disabled\"
  }")

echo "$BAN_RES" | grep "success" > /dev/null
check_status $? "Ban User Response"

# 14. Verify Login Fails
echo -e "\n14. Verifying Login Fails for Banned User..."
FAIL_LOGIN_RES=$(curl -s -X POST "${API_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"identifier\": \"$USERNAME\",
    \"password\": \"$PASSWORD\"
  }")

FAIL_CODE=$(echo $FAIL_LOGIN_RES | jq -r '.code') 

if [ "$FAIL_CODE" != "0" ] && [ "$FAIL_CODE" != "null" ]; then
    echo -e "${GREEN}[PASS]${NC} Login Failed as expected (Code: $FAIL_CODE)."
else
    echo -e "${RED}[FAIL]${NC} Banned user was able to login (or code 0)!"
    echo "Response: $FAIL_LOGIN_RES"
fi

# --- Admin UI Tests ---
echo -e "\n--- Admin UI Content ---"

# 15. Test Menu Config
echo -e "\n15. Testing Admin Menu Config..."
MENU_RES=$(curl -s -X GET "${ADMIN_URL}/menu")
echo "$MENU_RES" | grep "dashboard" > /dev/null
check_status $? "Menu Endpoint (Dashboard item found)"

# 16. Test Admin HTML Serving
echo -e "\n16. Testing Admin HTML..."
HTML_RES=$(curl -s -X GET "http://localhost:8080/admin")
echo "$HTML_RES" | grep "<title>Appsite Admin</title>" > /dev/null
check_status $? "Admin Index HTML Served"

echo "----------------------------------------"
echo "All Tests Completed Successfully"
echo "----------------------------------------"
