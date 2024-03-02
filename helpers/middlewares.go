package helpers

import (
	"log"
	"net/http"
	"points_mgmt/users"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("Token")

		if valid, err := IsValidToken(token); !valid {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return

		} else if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return

		}

		u, err := users.GetUserByToken(token, con)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.Set("idOrg", u.IdOrg)
		c.Next()
	}
}

func AuthRolePermission(role string) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("Token")

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passe um token válido no Header da requisição"})
			c.Abort()
			return
		}

		user, err := users.GetUserByToken(token, con)
		if err != nil {
			log.Println("Erro ao obter usuário")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		if user.Role != role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário sem permissão"})
			c.Abort()
			return
		}

		c.Next()
	}

}
