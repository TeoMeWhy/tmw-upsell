package api

import (
	"log"
	"net/http"
	"points_mgmt/customer"
	"points_mgmt/db"
	"points_mgmt/transaction"

	"github.com/gin-gonic/gin"
)

var con, _ = db.Connect()

func GetCustomer(c *gin.Context) {

	tx, err := con.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
		return
	}

	paramID := c.Query("id")
	if paramID != "" {
		customer, err := customer.GetCustomerByField(paramID, "UUID", tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if customer.UUID == "" {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusOK, customer)
		tx.Commit()
		return
	}

	paramCPF := c.Query("cpf")
	if paramCPF != "" {
		customer, err := customer.GetCustomerByField(paramCPF, "CPF", tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if customer.UUID == "" {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusOK, customer)
		tx.Commit()
		return
	}

	if len(c.Request.URL.Query()) > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilize parâmetros válidos"})
		return
	}

	customers, err := customer.GetCustomers(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, []customer.Customer{})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, customers)
}

func PostCustomer(c *gin.Context) {

	newCustomer := customer.Customer{}
	if err := c.BindJSON(&newCustomer); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	log.Println(newCustomer)

	tx, err := con.Begin()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
		return
	}

	oldCustomer := customer.Customer{}
	if newCustomer.CPF != "" {
		oldCustomer, err = customer.GetCustomerByField(newCustomer.CPF, "CPF", tx)
	} else if newCustomer.Email != "" {
		oldCustomer, err = customer.GetCustomerByField(newCustomer.Email, "Email", tx)
	} else if newCustomer.Name != "" {
		oldCustomer, err = customer.GetCustomerByField(newCustomer.Name, "Name", tx)
	}

	if err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível verificar se o usuário já existe"})
		return
	}

	if oldCustomer.UUID != "" {
		log.Println("Usuário já existente")
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário já existente"})
		return
	}

	if newCustomer, err = customer.CreateCustomer(newCustomer, tx); err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuário"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, newCustomer)
}

func PutAddUserPoints(c *gin.Context) {

	type Payload struct {
		UUID     string
		Points   int
		Products map[string]int
	}
	payload := Payload{}

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	if payload.UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário passar o parâmetro de id"})
		return
	}

	tx, err := con.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
		return
	}

	oldCustomer, err := customer.GetCustomerByField(payload.UUID, "UUID", tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if oldCustomer.UUID == "" {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não existente"})
		return
	}

	if payload.Points < 0 && (payload.Points*-1) > oldCustomer.Points {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pontos insuficientes"})
		return
	}

	points := payload.Points + oldCustomer.Points
	if err := customer.UpdateCustomerPoints(points, payload.UUID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	idTransaction, err := transaction.MakeTransaction(payload.Points, payload.UUID, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := transaction.MakeTransactionCart(idTransaction, payload.Products, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"points": payload.Points})
}
