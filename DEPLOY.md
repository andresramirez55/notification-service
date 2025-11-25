# Guía de Deploy - Notification Service

Este microservicio está diseñado para deployarse en un repositorio de GitHub separado y deployarse independientemente del backend principal.

## Pasos para Deploy

### 1. Crear Repositorio en GitHub

```bash
# En el directorio notification-service
cd notification-service
git init
git add .
git commit -m "Initial commit: Notification microservice with scheduler"
git branch -M main
git remote add origin https://github.com/tu-usuario/notification-service.git
git push -u origin main
```

### 2. Deploy en Railway

1. Ve a [Railway](https://railway.app)
2. Crea un nuevo proyecto
3. Selecciona "Deploy from GitHub repo"
4. Conecta el repositorio `notification-service`
5. Railway detectará automáticamente que es un proyecto Go

### 3. Configurar Variables de Entorno

En Railway, ve a la pestaña "Variables" y configura:

```env
DATABASE_URL=postgresql://user:password@host:port/dbname
SENDGRID_API_KEY=tu_api_key_de_sendgrid
FROM_EMAIL=noreply@tudominio.com
PORT=8081
```

**Importante**: 
- `DATABASE_URL` debe ser la misma base de datos que usa el backend principal
- Si usas PostgreSQL compartido, usa la misma URL de conexión

### 4. Verificar Deploy

Una vez deployado, verifica que el servicio esté funcionando:

```bash
curl https://tu-servicio.railway.app/health
```

Deberías recibir:
```json
{
  "status": "ok",
  "service": "notification-service",
  "version": "1.0.0",
  "time": "..."
}
```

## Funcionamiento del Scheduler

El scheduler se ejecuta automáticamente:

- **Cada hora**: Verifica eventos para hoy y mañana
- **A las 9:00 AM UTC**: Verificación especial para eventos del día siguiente
- **Eventos para mañana**: Envía notificaciones si `reminder_day_before = true`
- **Eventos para hoy**: Envía notificaciones si `reminder_day = true` (1 hora antes del evento)

## Monitoreo

Puedes monitorear los logs en Railway para ver:
- Cuándo se ejecuta el scheduler
- Qué eventos se encuentran
- Qué notificaciones se envían
- Errores si los hay

## Troubleshooting

### El scheduler no está funcionando

1. Verifica los logs en Railway
2. Asegúrate de que `DATABASE_URL` esté correctamente configurada
3. Verifica que la base de datos tenga eventos con `reminder_day` o `reminder_day_before = true`

### No se envían emails

1. Verifica que `SENDGRID_API_KEY` esté configurada
2. Verifica que `FROM_EMAIL` sea un email verificado en SendGrid
3. Revisa los logs para ver errores específicos

### No se conecta a la base de datos

1. Verifica que `DATABASE_URL` sea correcta
2. Asegúrate de que la base de datos permita conexiones desde Railway
3. Si usas PostgreSQL, verifica que el host sea accesible públicamente

