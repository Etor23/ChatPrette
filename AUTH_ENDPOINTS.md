# Configuración de Autenticación

Este backend usa Firebase Admin SDK para crear/verificar usuarios y emite JWTs propios para autenticación.

## Variables de Ambiente Necesarias

```env
# Base de datos MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=chat_app

# Server
PORT=8080

# Firebase credentials (ruta al archivo JSON)
FIREBASE_CREDENTIALS=./firebase-credentials.json

# JWT Secret Key (CAMBIAR en producción a una clave fuerte)
JWT_SECRET_KEY=your-super-secret-key-change-in-production
```

## Endpoints de Autenticación

### 1. POST /api/auth/register
**Público** - Crea un nuevo usuario

```json
Request:
{
  "email": "user@example.com",
  "password": "securePassword123",
  "username": "myusername",
  "avatar_url": "https://example.com/avatar.jpg" // opcional
}

Response (201 Created):
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "username": "myusername",
    "avatar_url": "https://example.com/avatar.jpg"
  },
  "is_new": true
}
```

### 2. POST /api/auth/login
**Público** - Inicia sesión con credenciales

```json
Request:
{
  "email": "user@example.com",
  "password": "securePassword123"
}

Response (200 OK):
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "username": "myusername",
    "avatar_url": "https://example.com/avatar.jpg"
  },
  "is_new": false
}
```

### 3. GET /api/auth/me
**Protegido** - Obtiene perfil del usuario autenticado

```
Headers:
Authorization: Bearer <token>

Response (200 OK):
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "username": "myusername",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

### 4. POST /api/auth/logout
**Protegido** - Cierra sesión (opcional)

```
Headers:
Authorization: Bearer <token>

Response (200 OK):
{
  "message": "Sesión cerrada correctamente"
}
```

### 5. POST /api/auth/refresh
**Protegido** - Refresca el token JWT

```
Headers:
Authorization: Bearer <token>

Response (200 OK):
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

## Códigos de Error

- `400 Bad Request` - Validación fallida (email inválido, password muy corto, etc.)
- `401 Unauthorized` - Token inválido, expirado o no proporcionado
- `404 Not Found` - Usuario no encontrado
- `409 Conflict` - Email o username ya en uso
- `500 Internal Server Error` - Error del servidor

## Flujo de Autenticación

1. **Usuario se registra**
   - POST /api/auth/register con email + password
   - Se crea en Firebase y MongoDB
   - Retorna JWT token de 24 horas

2. **Usuario inicia sesión**
   - POST /api/auth/login con email + password
   - Se verifica en Firebase
   - Retorna JWT token de 24 horas

3. **Usuario accede a rutas protegidas**
   - Envía Authorization: Bearer <token> en header
   - Middleware valida el token
   - Si es válido, permite acceso; sino, retorna 401

4. **Token expira (cada 24 horas)**
   - POST /api/auth/refresh para obtener nuevo token
   - Requiere Authorization header con token actual

## Seguridad

- Los tokens JWT expiran en **24 horas**
- Las contraseñas se almacenan en Firebase (no en MongoDB)
- Los emails y usernames son únicos
- CORS habilitado para localhost:3000 y localhost:5173

## Notas de Desarrollo

- El `JWT_SECRET_KEY` debe cambiarse en producción a una clave segura
- En producción, validar password en Firebase es más seguro usando REST API
- Considerar implementar refresh tokens y blacklist de tokens invalidados
- El middleware verifica tokens en: GET /auth/me, POST /auth/logout, POST /auth/refresh
