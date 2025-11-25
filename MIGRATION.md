# Migración a Microservicio de Notificaciones

Este documento explica cómo se migró la funcionalidad de notificaciones a un microservicio independiente.

## Cambios Realizados

### 1. Nuevo Microservicio (`notification-service/`)

- **Ubicación**: `/notification-service`
- **Puerto por defecto**: 8081
- **Funcionalidad**: Envío de emails usando SendGrid
- **API REST**: Endpoint `/api/v1/notifications/email`

### 2. Backend Principal (`backend-go/`)

- **Cambio**: Ya no envía emails directamente
- **Nuevo**: Usa un cliente HTTP para llamar al microservicio
- **Archivo nuevo**: `services/notification_client.go`
- **Modificado**: `services/notification_service.go` ahora usa el cliente

## Configuración

### Backend Principal

Agregar en `.env`:
```env
NOTIFICATION_SERVICE_URL=http://localhost:8081
```

Para producción:
```env
NOTIFICATION_SERVICE_URL=https://your-notification-service.railway.app
```

### Microservicio de Notificaciones

Crear `.env` en `notification-service/`:
```env
PORT=8081
SENDGRID_API_KEY=your_sendgrid_api_key_here
FROM_EMAIL=noreply@yourdomain.com
```

## Deploy

### Opción 1: Deploy Separado (Recomendado)

1. **Microservicio**: Deploy en Railway/Render con:
   - Variables: `SENDGRID_API_KEY`, `FROM_EMAIL`
   - Puerto: 8081 (o el que asigne la plataforma)

2. **Backend Principal**: Deploy con:
   - Variable: `NOTIFICATION_SERVICE_URL` apuntando al microservicio

### Opción 2: Docker Compose (Desarrollo)

Puedes usar docker-compose para levantar ambos servicios localmente.

## Ventajas

✅ **Separación de responsabilidades**: El backend principal no necesita conocer SendGrid
✅ **Escalabilidad**: El microservicio puede escalarse independientemente
✅ **Mantenibilidad**: Cambios en notificaciones no afectan el backend principal
✅ **Deploy independiente**: Puedes actualizar notificaciones sin tocar el backend

## Próximos Pasos

- [ ] Agregar soporte para WhatsApp en el microservicio (opcional)
- [ ] Implementar retry logic en el cliente
- [ ] Agregar métricas y logging
- [ ] Implementar rate limiting

