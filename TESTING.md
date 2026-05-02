# Testing Autenticación - Ejemplos de CURL

Este archivo contiene ejemplos para probar los endpoints de autenticación usando curl desde PowerShell.

## Configuración Previa

Asegúrate que:
1. MongoDB está corriendo: `localhost:27017`
2. El backend está activo: `go run .` en http://localhost:8080
3. Tienes `curl` instalado (disponible por defecto en Windows 10+)

## 1. Registrar Nuevo Usuario

```powershell
$body = @{
    email = "testuser@example.com"
    password = "TestPassword123"
    username = "testuser"
    avatar_url = ""
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/api/auth/register" `
  -Method POST `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body
```

**Respuesta esperada (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "testuser@example.com",
    "username": "testuser",
    "avatar_url": ""
  },
  "is_new": true
}
```

**Posibles errores:**
- `400` - Email inválido, password muy corto, username no proporcionado
- `409` - Email o username ya existe

---

## 2. Login

```powershell
$body = @{
    email = "testuser@example.com"
    password = "TestPassword123"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/api/auth/login" `
  -Method POST `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "testuser@example.com",
    "username": "testuser",
    "avatar_url": ""
  },
  "is_new": false
}
```

**Posibles errores:**
- `401` - Email o password incorrecto
- `404` - Usuario no registrado

---

## 3. Obtener Perfil (GET /auth/me)

```powershell
$token = "eyJhbGciOiJIUzI1NiIs..." # Reemplazar con token del login/register

Invoke-WebRequest -Uri "http://localhost:8080/api/auth/me" `
  -Method GET `
  -Headers @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer $token"
  }
```

**Respuesta esperada (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "testuser@example.com",
  "username": "testuser",
  "avatar_url": ""
}
```

**Posibles errores:**
- `401` - Token no proporcionado o inválido
- `404` - Usuario no encontrado

---

## 4. Refrescar Token (POST /auth/refresh)

```powershell
$token = "eyJhbGciOiJIUzI1NiIs..." # Token actual

Invoke-WebRequest -Uri "http://localhost:8080/api/auth/refresh" `
  -Method POST `
  -Headers @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer $token"
  }
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR...", # Nuevo token
  "expires_in": 86400
}
```

---

## 5. Logout (POST /auth/logout)

```powershell
$token = "eyJhbGciOiJIUzI1NiIs..." # Token actual

Invoke-WebRequest -Uri "http://localhost:8080/api/auth/logout" `
  -Method POST `
  -Headers @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer $token"
  }
```

**Respuesta esperada (200 OK):**
```json
{
  "message": "Sesión cerrada correctamente"
}
```

---

## Verificar MongoDB

Acceder a la BD y ver los usuarios creados:

```bash
mongosh mongodb://localhost:27017
> use chat_app
> db.users.find().pretty()
```

---

## Notas

- Los tokens expiran en **24 horas**
- Todos los campos JSON son sensibles a mayúsculas/minúsculas
- Las contraseñas deben tener al menos 6 caracteres
- Los emails deben ser válidos
- Los usernames deben tener entre 3 y 32 caracteres
