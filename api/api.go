package api

import (
	"log"
	"net/http"
	"points_mgmt/customers"
	"points_mgmt/db"
	"points_mgmt/transaction"

	"github.com/gin-gonic/gin"
)

var con, _ = db.Connect()

func GetCustomer(c *gin.Context) {

	idOrg, ok := c.Get("idOrg")
	if !ok {
		log.Println("idOrg nao foi passada")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "idOrg nao foi passada"})
		return
	}

	mapFilters := map[string]string{"IdOrg": idOrg.(string)}

	tx, err := con.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
		return
	}

	paramID := c.Query("id")
	if paramID != "" {

		mapFilters["UUID"] = paramID

		customer, err := customers.GetCustomerByField(mapFilters, tx)
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

		mapFilters["CPF"] = paramCPF
		customer, err := customers.GetCustomerByField(mapFilters, tx)
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

	nCustomers, err := customers.GetCustomers(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, []customers.Customer{})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, nCustomers)
}

func PostCustomer(c *gin.Context) {

	idOrg, ok := c.Get("idOrg")
	if !ok {
		log.Println("idOrg nao foi passada")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "idOrg nao foi passada"})
		return
	}
	mapFilters := map[string]string{"IdOrg": idOrg.(string)}

	newCustomer := customers.Customer{}
	if err := c.BindJSON(&newCustomer); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	newCustomer.IdOrg = idOrg.(string)

	tx, err := con.Begin()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
		return
	}

	oldCustomer := customers.Customer{}
	if newCustomer.CPF != "" {
		mapFilters["CPF"] = newCustomer.CPF
	} else if newCustomer.Email != "" {
		mapFilters["Email"] = newCustomer.Email
	} else if newCustomer.Name != "" {
		mapFilters["Name"] = newCustomer.Name
	}

	oldCustomer, err = customers.GetCustomerByField(mapFilters, tx)
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

	if newCustomer, err = customers.CreateCustomer(newCustomer, tx); err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o usuário"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, newCustomer)
}

func DeleteCustomer(c *gin.Context) {

	idOrg, ok := c.Get("idOrg")
	if !ok {
		log.Println("idOrg nao foi passada")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "idOrg nao foi passada"})
		return
	}
	mapFilters := map[string]string{"IdOrg": idOrg.(string)}

	paramID := c.Query("id")
	if paramID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entre com um usuário válido"})
		return
	}

	mapFilters["UUID"] = paramID

	tx, err := con.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível iniciar a transaction"})
		return
	}

	statement, err := tx.Prepare("DELETE FROM tb_customers WHERE UUID = ? AND IdOrg = ?")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível iniciar o statement"})
		return
	}

	res, err := statement.Exec(mapFilters["UUID"], mapFilters["IdOrg"])
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível executar o statement"})
		return
	}

	row, _ := res.RowsAffected()
	if row == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"sucess": "Usuário removido com sucesso"})
}

func PutCustomerEmail(c *gin.Context) {

	idOrg, ok := c.Get("idOrg")
	if !ok {
		log.Println("idOrg nao foi passada")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "idOrg nao foi passada"})
		return
	}

	type EmailPayload struct {
		UUID  string
		Email string
	}
	payload := EmailPayload{}

	err := c.Bind(&payload)
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

	// if payload.Email != "" {
	// 	oldCustomer, err := customer.GetCustomerByField(payload.Email, "Email", tx)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 		return
	// 	}

	// 	if oldCustomer.UUID != "" {
	// 		tx.Rollback()
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "email já cadastrado por outro usuário"})
	// 		return
	// 	}
	// }

	customer := customers.Customer{
		UUID:  payload.UUID,
		Email: payload.Email,
		IdOrg: idOrg.(string),
	}

	if err := customers.UpdateCustomerEmail(customer, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o Email"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, nil)
}

func PutCustomer(c *gin.Context) {
	payload := customers.Customer{}

	err := c.Bind(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter os valores do body"})
		return
	}

	tx, err := con.Begin()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao Erro ao abrir a transaction de PutCustomer"})
		return
	}

	if err := customers.UpdateCustomer(payload, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	tx.Commit()
}

func PutAddCustomerPoints(c *gin.Context) {

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

	oldCustomer, err := customers.GetCustomerByField(map[string]string{"UUID": payload.UUID}, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	log.Println(oldCustomer)

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
	if err := customers.UpdateCustomerPoints(points, payload.UUID, tx); err != nil {
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

func GetCustomerTransactions(c *gin.Context) {

	idCustomer := c.Query("id")
	if idCustomer == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entre com usuário válido"})
		return
	}

	log.Println(idCustomer)

	userTransactions, err := transaction.GetCustomerTransactions(idCustomer, con)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	log.Println(userTransactions)

	c.JSON(http.StatusOK, gin.H{"transactions": userTransactions})

}
