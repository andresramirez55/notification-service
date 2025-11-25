# Setup GitHub Repository

El repositorio Git ya está inicializado en este directorio. Sigue estos pasos para crear el repositorio en GitHub y hacer push:

## Pasos

### 1. Crear el repositorio en GitHub

1. Ve a https://github.com/new
2. Nombre del repositorio: `notification-service` (o el nombre que prefieras)
3. **NO** inicialices con README, .gitignore o licencia (ya los tenemos)
4. Haz clic en "Create repository"

### 2. Conectar y hacer push

Una vez creado el repositorio en GitHub, ejecuta estos comandos:

```bash
cd /Users/andresramirez/calendar-starter/notification-service

# Agregar el remote (reemplaza TU_USUARIO con tu usuario de GitHub)
git remote add origin https://github.com/TU_USUARIO/notification-service.git

# O si prefieres SSH:
# git remote add origin git@github.com:TU_USUARIO/notification-service.git

# Hacer push
git branch -M main
git push -u origin main
```

### 3. Verificar

Después del push, verifica que todo esté bien:

```bash
git remote -v
git log --oneline
```

## Siguiente paso: Deploy en Railway

Una vez que el código esté en GitHub:

1. Ve a [Railway](https://railway.app)
2. Crea un nuevo proyecto
3. Selecciona "Deploy from GitHub repo"
4. Conecta el repositorio `notification-service`
5. Configura las variables de entorno (ver DEPLOY.md)

¡Listo! El microservicio estará deployado y funcionando automáticamente.

