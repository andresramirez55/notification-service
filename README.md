# Notification Service

Microservicio independiente para el envío de notificaciones por email con scheduler automático.

## Características

- ✅ Envío de emails usando SendGrid
- 🔄 API REST simple y clara
- ⏰ Scheduler automático que busca eventos para hoy y mañana
- 📅 Notificaciones para eventos del día siguiente (mañana)
- 📅 Notificaciones para eventos del mismo día (hoy)
- 🐳 Dockerizado y listo para deploy
- 🏥 Health check endpoint

## Configuración

### Variables de Entorno

Crea un archivo `.env` o configura las siguientes variables:

```env
PORT=8081
SENDGRID_API_KEY=your_sendgrid_api_key_here
FROM_EMAIL=noreply@calendar.com
DATABASE_URL=postgresql://user:password@host:port/dbname
# O para desarrollo local con SQLite:
# DATABASE_URL=calendar.db
```

**Importante**: El microservicio necesita acceso a la misma base de datos que el backend principal para consultar eventos.

## Uso

### Desarrollo Local

```bash
# Instalar dependencias
go mod download

# Ejecutar
go run main.go
```

### Con Docker

```bash
# Build
docker build -t notification-service .

# Run
docker run -p 8081:8081 --env-file .env notification-service
```

## API Endpoints

### Health Check
```
GET /health
```

### Enviar Email
```
POST /api/v1/notifications/email
Content-Type: application/json

{
  "to": "user@example.com",
  "subject": "Recordatorio: Evento mañana",
  "body": "Hola! Te recordamos que mañana tenés...",
  "recipient_name": "Usuario" // opcional
}
```

### Status
```
GET /api/v1/notifications/status
```

## Funcionamiento del Scheduler

El scheduler se ejecuta automáticamente y:

1. **Cada hora**: Verifica eventos para hoy y mañana
2. **A las 9:00 AM UTC**: Verificación especial para eventos del día siguiente
3. **Eventos para mañana**: Envía notificaciones si `reminder_day_before = true`
4. **Eventos para hoy**: Envía notificaciones si `reminder_day = true` (1 hora antes del evento)

## Deploy

Este servicio está diseñado para deployarse independientemente del backend principal. Puedes deployarlo en:

- Railway
- Render
- Heroku
- Cualquier plataforma que soporte Docker/Go

### Variables de Entorno para Deploy

Asegúrate de configurar:
- `DATABASE_URL`: URL de la base de datos (debe ser la misma que usa el backend principal)
- `SENDGRID_API_KEY`: API key de SendGrid
- `FROM_EMAIL`: Email desde el cual se envían las notificaciones
- `PORT`: Puerto del servicio (opcional, por defecto 8081)

### Deploy en Railway

1. Crea un nuevo proyecto en Railway
2. Conecta este repositorio
3. Configura las variables de entorno
4. El servicio se desplegará automáticamente

**Nota**: Este microservicio debe tener acceso a la misma base de datos que el backend principal para poder consultar los eventos.

