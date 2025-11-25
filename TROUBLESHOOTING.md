# Troubleshooting - Notification Service

## Error: "Service unavailable" en Railway

Si ves este error durante el deploy, aquí están las causas más comunes y soluciones:

### 1. Variable `DATABASE_URL` no configurada

**Síntoma**: El servicio crashea al iniciar con error de conexión a base de datos.

**Solución**:
1. Ve a Railway → Tu proyecto → Variables
2. Agrega la variable `DATABASE_URL` con la URL de tu base de datos PostgreSQL
3. Formato: `postgresql://user:password@host:port/dbname`

**Nota**: Debe ser la misma base de datos que usa el backend principal.

### 2. Base de datos no accesible

**Síntoma**: El servicio inicia pero el health check retorna `database: disconnected`.

**Solución**:
- Verifica que la base de datos esté corriendo
- Verifica que la URL de conexión sea correcta
- Si usas Railway PostgreSQL, asegúrate de que el servicio tenga acceso a la misma base de datos

### 3. Puerto incorrecto

**Síntoma**: El servicio no responde a los health checks.

**Solución**:
- Railway asigna el puerto automáticamente via la variable `PORT`
- No necesitas configurarlo manualmente
- El servicio usa `os.Getenv("PORT")` que Railway inyecta automáticamente

### 4. Build falla

**Síntoma**: El build no completa.

**Solución**:
- Verifica que `go.mod` y `go.sum` estén en el repositorio
- Railway debería detectar automáticamente que es un proyecto Go
- Si no, verifica que el `Dockerfile` esté presente

## Verificar el estado del servicio

### Health Check

```bash
curl https://tu-servicio.railway.app/health
```

Respuesta esperada:
```json
{
  "status": "ok",
  "service": "notification-service",
  "version": "1.0.0",
  "database": "connected",
  "time": "2024-11-25T..."
}
```

Si `database` es `not_initialized` o `disconnected`, el problema es la conexión a la base de datos.

### Ver logs en Railway

1. Ve a Railway → Tu proyecto → Deployments
2. Haz clic en el deployment más reciente
3. Ve a la pestaña "Logs"
4. Busca errores relacionados con:
   - `Failed to connect to database`
   - `DATABASE_URL`
   - `PostgreSQL`

## Variables de entorno requeridas

Asegúrate de tener estas variables configuradas en Railway:

```env
DATABASE_URL=postgresql://user:password@host:port/dbname  # REQUERIDO
SENDGRID_API_KEY=tu_api_key                                # REQUERIDO para enviar emails
FROM_EMAIL=noreply@tudominio.com                           # REQUERIDO para enviar emails
PORT=8081                                                   # Opcional, Railway lo asigna automáticamente
```

## El servicio inicia pero el scheduler no funciona

**Causa**: La base de datos no está conectada.

**Solución**:
1. Verifica que `DATABASE_URL` esté configurada correctamente
2. Verifica los logs para ver si hay errores de conexión
3. El scheduler solo se inicia si la conexión a la BD es exitosa

## Verificar que el scheduler está corriendo

En los logs deberías ver:
```
🚀 Starting notification scheduler...
✅ Notification scheduler started
```

Si no ves estos mensajes, el scheduler no se inició (probablemente por falta de conexión a BD).

## Próximos pasos si el problema persiste

1. Verifica los logs completos en Railway
2. Prueba el health check endpoint
3. Verifica que todas las variables de entorno estén configuradas
4. Asegúrate de que la base de datos esté accesible desde Railway

