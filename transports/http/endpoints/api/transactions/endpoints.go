package transactions

import (
	"bux-wallet/domain"
	"bux-wallet/domain/transactions"
	"bux-wallet/domain/users"
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
		user.GET("", h.getTransactions)
		user.GET("/:id", h.getTransaction)
	}
}

// Get all user transactions.
//
//	@Summary Get all transactions.
//	@Tags transaction
//	@Produce json
//	@Success 200 {object} []buxclient.Transaction
//	@Router /api/v1/transaction [get]
func (h *handler) getTransactions(c *gin.Context) {
	// Get user transactions.
	transactions, err := h.tService.GetTransactions(c.GetString(auth.SessionAccessKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Get specific transactions.
//
//	@Summary Get transaction by id.
//	@Tags transaction
//	@Produce json
//	@Success 200 {object} buxclient.FullTransaction
//	@Router /api/v1/transaction/{id} [get]
//	@Param id path string true "Transaction id"
func (h *handler) getTransaction(c *gin.Context) {
	transactionId := c.Param("id")

	// Get transaction by id.
	transaction, err := h.tService.GetTransaction(c.GetString(auth.SessionAccessKey), transactionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	c.JSON(http.StatusOK, transaction)
}