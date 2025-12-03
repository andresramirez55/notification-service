# Cómo Probar el Servicio de Notificaciones

## Verificar que el Evento Esté Configurado Correctamente

Antes de probar, asegúrate de que tu evento tenga:

1. **Fecha de hoy** (para notificaciones del mismo día)
2. **`reminder_day = true`** (para recibir notificación el día del evento)
3. **`reminder_day_before = true`** (para recibir notificación el día anterior)
4. **Email válido** configurado en el evento

## Métodos de Prueba

### 1. Verificar Health del Servicio

```bash
curl https://tu-notification-service.railway.app/health
```

Deberías recibir:
```json
{
  "status": "ok",
  "service": "notification-service",
  "database": "connected",
  ...
}
```

### 2. Forzar Verificación Manual (Recomendado)

El servicio ahora tiene un endpoint para forzar la verificación inmediatamente:

```bash
curl -X POST https://tu-notification-service.railway.app/api/v1/notifications/check
```

Esto:
- Busca eventos para **hoy** con `reminder_day = true`
- Busca eventos para **mañana** con `reminder_day_before = true`
- Envía emails automáticamente si encuentra eventos

### 3. Ver Logs del Servicio

En Railway, ve a:
1. Tu proyecto → Deployments
2. Haz clic en el deployment más reciente
3. Ve a la pestaña "Logs"

Deberías ver:
```
🔍 Checking for notifications to send...
📅 Checking for day-before reminders (tomorrow's events)...
📅 Checking for same-day reminders (today's events)...
📧 Found X events for today, sending notifications...
✅ Email sent successfully to...
```

### 4. Enviar Email de Prueba Directo

Para probar que SendGrid funciona:

```bash
curl -X POST https://tu-notification-service.railway.app/api/v1/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "tu-email@ejemplo.com",
    "subject": "Test de Notificación",
    "body": "Este es un email de prueba del servicio de notificaciones",
    "recipient_name": "Tu Nombre"
  }'
```

## Verificar que el Evento Sea Detectado

El scheduler busca eventos que cumplan estas condiciones:

**Para eventos de hoy:**
- `date = hoy`
- `reminder_day = true`
- Se envía 1 hora antes del evento (o en la mañana si es evento de todo el día)

**Para eventos de mañana:**
- `date = mañana`
- `reminder_day_before = true`
- Se envía automáticamente

## Troubleshooting

### No se encuentran eventos

1. Verifica que el evento tenga la fecha correcta
2. Verifica que `reminder_day` o `reminder_day_before` estén en `true`
3. Verifica los logs del servicio para ver qué está buscando

### No se envía el email

1. Verifica que `SENDGRID_API_KEY` esté configurada en Railway
2. Verifica que `FROM_EMAIL` sea un email verificado en SendGrid
3. Revisa los logs para ver errores específicos

### El endpoint /check no funciona

1. Verifica que la base de datos esté conectada (health check debe mostrar `"database": "connected"`)
2. Verifica los logs para ver si hay errores de conexión

## Script de Prueba

Usa el script incluido:

```bash
cd notification-service
./test_notifications.sh https://tu-notification-service.railway.app
```

O para local:

```bash
./test_notifications.sh http://localhost:8081
```

