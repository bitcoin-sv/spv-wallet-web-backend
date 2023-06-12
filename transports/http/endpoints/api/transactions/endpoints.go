package transactions

import (
	"bux-wallet/domain"
	"bux-wallet/domain/transactions"
	"bux-wallet/domain/users"
	"fmt"
	"net/http"

	"bux-wallet/transports/http/auth"
	"bux-wallet/transports/http/endpoints/api"
	router "bux-wallet/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

type handler struct {
	uService users.UserService
	tService transactions.TransactionService
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services) router.ApiEndpoints {
	return &handler{
		uService: *s.UsersService,
		tService: *s.TransactionsService,
	}
}

// RegisterApiEndpoints registers routes that are part of service API.
func (h *handler) RegisterApiEndpoints(router *gin.RouterGroup) {
	user := router.Group("/transaction")
	{
		user.POST("", h.createTransaction)
		user.GET("", h.getTransactions)
		user.GET("/:id", h.getTransaction)
	}
}

func (h *handler) createTransaction(c *gin.Context) {
	var reqTransaction CreateTransaction
	err := c.Bind(&reqTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	// Validate user.
	xpriv, err := h.uService.GetUserXpriv(c.GetInt("userId"), reqTransaction.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	// Create transaction.
	transaction, err := h.tService.CreateTransaction(xpriv, reqTransaction.Recipient, reqTransaction.Satoshis)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	fmt.Println("Transaction: ", transaction)
}

func (h *handler) getTransactions(c *gin.Context) {
	// Get user transactions.
	transactions, err := h.tService.GetTransactions(c.GetString(auth.SessionAccessKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *handler) getTransaction(c *gin.Context) {
	transactionId := c.Param("id")

	// Get user transactions.
	transactions, err := h.tService.GetTransaction(c.GetString(auth.SessionAccessKey), transactionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	c.JSON(http.StatusOK, transactions)
}
