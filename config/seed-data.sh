#!/bin/bash

# Aegis Test Data Seeding Script
# This script populates the database with sample users, roles, and permissions

set -e

API_BASE="${API_BASE:-http://localhost:3100/api/aegis}"

echo "🌱 Seeding Aegis with test data..."
echo "📡 API Base URL: $API_BASE"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create Permissions
echo -e "${BLUE}Creating Permissions...${NC}"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"read:users","description":"Can read user information"}' > /dev/null
echo -e "${GREEN}✓${NC} read:users"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"write:users","description":"Can create and update users"}' > /dev/null
echo -e "${GREEN}✓${NC} write:users"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"delete:users","description":"Can delete users"}' > /dev/null
echo -e "${GREEN}✓${NC} delete:users"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"read:reports","description":"Can view reports"}' > /dev/null
echo -e "${GREEN}✓${NC} read:reports"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"write:reports","description":"Can create and edit reports"}' > /dev/null
echo -e "${GREEN}✓${NC} write:reports"

curl -s -X POST "$API_BASE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"name":"manage:system","description":"Can manage system settings"}' > /dev/null
echo -e "${GREEN}✓${NC} manage:system"

echo ""

# Create Roles
echo -e "${BLUE}Creating Roles...${NC}"

curl -s -X POST "$API_BASE/roles" \
  -H "Content-Type: application/json" \
  -d '{"name":"admin","description":"System administrator with full access"}' > /dev/null
echo -e "${GREEN}✓${NC} admin"

curl -s -X POST "$API_BASE/roles" \
  -H "Content-Type: application/json" \
  -d '{"name":"manager","description":"Manager with user management access"}' > /dev/null
echo -e "${GREEN}✓${NC} manager"

curl -s -X POST "$API_BASE/roles" \
  -H "Content-Type: application/json" \
  -d '{"name":"viewer","description":"Read-only access to view data"}' > /dev/null
echo -e "${GREEN}✓${NC} viewer"

curl -s -X POST "$API_BASE/roles" \
  -H "Content-Type: application/json" \
  -d '{"name":"analyst","description":"Can view and create reports"}' > /dev/null
echo -e "${GREEN}✓${NC} analyst"

echo ""

# Create Users and Assign Roles/Permissions
echo -e "${BLUE}Creating Users...${NC}"

# Alice - Admin
ALICE=$(curl -s -X POST "$API_BASE/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "subject":"alice@aegis.com",
    "password":"Password123!",
    "additionalInfo":"{\"firstName\":\"Alice\",\"lastName\":\"Anderson\",\"department\":\"IT\",\"title\":\"System Administrator\"}"
  }' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST "$API_BASE/users/$ALICE/roles" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}' > /dev/null

curl -s -X POST "$API_BASE/users/$ALICE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"manage:system"}' > /dev/null

echo -e "${GREEN}✓${NC} alice@aegis.com (Admin)"

# Bob - Manager
BOB=$(curl -s -X POST "$API_BASE/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "subject":"bob@aegis.com",
    "password":"Password123!",
    "additionalInfo":"{\"firstName\":\"Bob\",\"lastName\":\"Brown\",\"department\":\"Sales\",\"title\":\"Sales Manager\"}"
  }' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST "$API_BASE/users/$BOB/roles" \
  -H "Content-Type: application/json" \
  -d '{"role":"manager"}' > /dev/null

curl -s -X POST "$API_BASE/users/$BOB/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"read:users"}' > /dev/null

curl -s -X POST "$API_BASE/users/$BOB/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"write:users"}' > /dev/null

echo -e "${GREEN}✓${NC} bob@aegis.com (Manager)"

# Carol - Viewer
CAROL=$(curl -s -X POST "$API_BASE/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "subject":"carol@aegis.com",
    "password":"Password123!",
    "additionalInfo":"{\"firstName\":\"Carol\",\"lastName\":\"Chen\",\"department\":\"Marketing\",\"title\":\"Marketing Specialist\"}"
  }' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST "$API_BASE/users/$CAROL/roles" \
  -H "Content-Type: application/json" \
  -d '{"role":"viewer"}' > /dev/null

curl -s -X POST "$API_BASE/users/$CAROL/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"read:reports"}' > /dev/null

echo -e "${GREEN}✓${NC} carol@aegis.com (Viewer)"

# David - Viewer
DAVID=$(curl -s -X POST "$API_BASE/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "subject":"david@aegis.com",
    "password":"Password123!",
    "additionalInfo":"{\"firstName\":\"David\",\"lastName\":\"Davis\",\"department\":\"Engineering\",\"title\":\"Software Engineer\"}"
  }' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST "$API_BASE/users/$DAVID/roles" \
  -H "Content-Type: application/json" \
  -d '{"role":"viewer"}' > /dev/null

curl -s -X POST "$API_BASE/users/$DAVID/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"read:users"}' > /dev/null

echo -e "${GREEN}✓${NC} david@aegis.com (Viewer)"

# Eve - Analyst
EVE=$(curl -s -X POST "$API_BASE/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "subject":"eve@aegis.com",
    "password":"Password123!",
    "additionalInfo":"{\"firstName\":\"Eve\",\"lastName\":\"Evans\",\"department\":\"Analytics\",\"title\":\"Data Analyst\"}"
  }' | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST "$API_BASE/users/$EVE/roles" \
  -H "Content-Type: application/json" \
  -d '{"role":"analyst"}' > /dev/null

curl -s -X POST "$API_BASE/users/$EVE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"read:reports"}' > /dev/null

curl -s -X POST "$API_BASE/users/$EVE/permissions" \
  -H "Content-Type: application/json" \
  -d '{"permission":"write:reports"}' > /dev/null

echo -e "${GREEN}✓${NC} eve@aegis.com (Analyst)"

echo ""
echo -e "${GREEN}✅ Test data seeded successfully!${NC}"
echo ""
echo "Summary:"
echo "  • 6 Permissions created"
echo "  • 4 Roles created"
echo "  • 5 Users created with role/permission assignments"
echo ""
echo "All users have password: Password123!"
echo ""
echo "🌐 Open http://localhost to view the UI"
