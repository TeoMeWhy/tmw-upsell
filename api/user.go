package api

import (
	"log"
	"net/http"
	"points_mgmt/users"

	"github.com/gin-gonic/gin"
)

func PostUsers(c *gin.Context) {

	user := &users.User{}

	if err := c.Bind(&user); err != nil {
		log.Println("Erro ao realizar o parser do body", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := users.NewUser(user.Name, user.Email, user.IdOrg, user.Role)
	if err != nil {
		log.Println("Erro ao criar novo usuário", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	tx, err := con.Begin()
	if err != nil {
		log.Println("Erro ao abrir a transação de novo usuário", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	exists, err := users.UserExists(user.Email, user.IdOrg, tx)
	if err != nil {
		tx.Rollback()
		log.Println("Erro ao buscar usuário", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if exists {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário ja existente"})
		return
	}

	if err := user.CreateUser(tx); err != nil {
		log.Println("Erro ao criar usuário", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"success": "true"})

}
