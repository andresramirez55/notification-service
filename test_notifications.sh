#!/bin/bash

# Script para probar el servicio de notificaciones
# Uso: ./test_notifications.sh [URL_DEL_SERVICIO]

NOTIFICATION_SERVICE_URL="${1:-http://localhost:8081}"

echo "🧪 Testing Notification Service"
echo "================================"
echo "Service URL: $NOTIFICATION_SERVICE_URL"
echo ""

# Test 1: Health check
echo "1. Testing health endpoint..."
curl -s "$NOTIFICATION_SERVICE_URL/health" | jq . 2>/dev/null || curl -s "$NOTIFICATION_SERVICE_URL/health"
echo ""
echo ""

# Test 2: Status
echo "2. Testing status endpoint..."
curl -s "$NOTIFICATION_SERVICE_URL/api/v1/notifications/status" | jq . 2>/dev/null || curl -s "$NOTIFICATION_SERVICE_URL/api/v1/notifications/status"
echo ""
echo ""

# Test 3: Force notification check (busca eventos para hoy y mañana)
echo "3. Forcing notification check (busca eventos para hoy y mañana)..."
curl -X POST -s "$NOTIFICATION_SERVICE_URL/api/v1/notifications/check" | jq . 2>/dev/null || curl -X POST -s "$NOTIFICATION_SERVICE_URL/api/v1/notifications/check"
echo ""
echo ""

# Test 4: Send test email (opcional)
echo "4. To send a test email, use:"
echo "   curl -X POST $NOTIFICATION_SERVICE_URL/api/v1/notifications/email \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"to\":\"tu-email@ejemplo.com\",\"subject\":\"Test\",\"body\":\"Test email\"}'"
echo ""

echo "✅ Tests completed!"
echo ""
echo "📋 Next steps:"
echo "   1. Verifica los logs del servicio para ver qué eventos encontró"
echo "   2. Revisa tu email para ver si recibiste la notificación"
echo "   3. Asegúrate de que el evento tenga reminder_day=true o reminder_day_before=true"

