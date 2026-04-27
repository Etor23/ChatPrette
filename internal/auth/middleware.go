package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextFirebaseUID = "firebase_uid"
	ContextEmail       = "email"
)

// Middleware verifica el token de Firebase en cada request
func Middleware(provider *FirebaseProvider) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Obtener header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Falta el header Authorization",
			})
			return
		}

		// 2. Extraer token: "Bearer eyJhb..."
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato inválido. Usa: Bearer <token>",
			})
			return
		}

		// 3. Verificar con Firebase
		tokenInfo, err := provider.VerifyToken(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido o expirado",
			})
			return
		}

		// 4. Guardar en contexto para que los handlers lo usen
		c.Set(ContextFirebaseUID, tokenInfo.UID)
		c.Set(ContextEmail, tokenInfo.Email)

		c.Next()
	}
}

// Helpers para extraer datos del contexto en los handlers
func GetFirebaseUID(c *gin.Context) string {
	uid, _ := c.Get(ContextFirebaseUID)
	return uid.(string)
}

func GetEmail(c *gin.Context) string {
	email, _ := c.Get(ContextEmail)
	return email.(string)
}
