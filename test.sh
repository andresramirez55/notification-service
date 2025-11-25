#!/bin/bash

# Test script for notification service

echo "🧪 Testing Notification Service"
echo "================================"

# Test health endpoint
echo ""
echo "1. Testing health endpoint..."
curl -s http://localhost:8081/health | jq .
echo ""

# Test email endpoint
echo "2. Testing email endpoint..."
curl -X POST http://localhost:8081/api/v1/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Test Email",
    "body": "This is a test email from the notification service",
    "recipient_name": "Test User"
  }' | jq .

echo ""
echo "✅ Tests completed!"

