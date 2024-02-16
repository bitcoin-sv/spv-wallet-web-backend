package transactions

import (
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/transactions"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/notification"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/websocket"

	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/rs/zerolog"

	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/auth"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

type handler struct {
	uService users.UserService
	tService transactions.TransactionService
	log      *zerolog.Logger
	ws       websocket.Server
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services, log *zerolog.Logger, ws websocket.Server) router.ApiEndpoints {
	return &handler{
		uService: *s.UsersService,
		tService: *s.TransactionsService,
		log:      log,
		ws:       ws,
	}
}

// RegisterApiEndpoints registers routes that are part of service API.
func (h *handler) RegisterApiEndpoints(router *gin.RouterGroup) {
	user := router.Group("/transaction")
	{
		user.GET("", h.getTransactions)
		user.POST("", h.createTransaction)
		user.GET("/:id", h.getTransaction)
	}
}

// Get all user transactions.
//
//	@Summary Get all transactions.
//	@Tags transaction
//	@Produce json
//	@Success 200 {object} []spvwallet.Transaction
//	@Router /api/v1/transaction [get]
func (h *handler) getTransactions(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	orderBy := c.Query("order")
	sort := c.Query("sort")

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		pageNumber = 1
	}

	pageSizeNumber, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeNumber = 10
	}

	queryParam := transports.QueryParams{
		Page:          pageNumber,
		PageSize:      pageSizeNumber,
		OrderByField:  orderBy,
		SortDirection: sort,
	}

	// Get user transactions.
	txs, err := h.tService.GetTransactions(c.GetString(auth.SessionAccessKey), c.GetString(auth.SessionUserPaymail), queryParam)
	if err != nil {
		h.log.Error().Msgf("An error occurred while trying to get a list of transactions: %s", err)
		c.JSON(http.StatusInternalServerError, "An error occurred while trying to get a list of transactions")
		return
	}

	c.JSON(http.StatusOK, txs)
}

// Get specific transactions.
//
//	@Summary Get transaction by id.
//	@Tags transaction
//	@Produce json
//	@Success 200 {object} spvwallet.FullTransaction
//	@Router /api/v1/transaction/{id} [get]
//	@Param id path string true "Transaction id"
func (h *handler) getTransaction(c *gin.Context) {
	transactionId := c.Param("id")

	// Get transaction by id.
	transaction, err := h.tService.GetTransaction(c.GetString(auth.SessionAccessKey), transactionId, c.GetString(auth.SessionUserPaymail))
	if err != nil {
		h.log.Error().Msgf("An error occurred while trying to get transaction details: %s", err)
		c.JSON(http.StatusInternalServerError, "An error occurred while trying to get transaction details")
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Create transactions.
//
//	@Summary Create transaction.
//	@Tags transaction
//	@Produce json
//	@Success 200 {object} spvwallet.FullTransaction
//	@Router /api/v1/transaction [post]
//	@Param data body CreateTransaction true "Create transaction data"
func (h *handler) createTransaction(c *gin.Context) {
	var reqTransaction CreateTransaction
	err := c.Bind(&reqTransaction)
	if err != nil {
		h.log.Error().Msgf("Invalid payload: %s", err)
		c.JSON(http.StatusBadRequest, "Invalid request. Please check transaction details")
		return
	}

	// Validate user.
	xpriv, err := h.uService.GetUserXpriv(c.GetInt(auth.SessionUserId), reqTransaction.Password)
	if err != nil {
		h.log.Error().Msgf("Invalid password: %s", err)
		c.JSON(http.StatusBadRequest, "Invalid password.")
		return
	}

	events := make(chan notification.TransactionEvent)
	err = h.tService.CreateTransaction(c.GetString(auth.SessionUserPaymail), xpriv, reqTransaction.Recipient, reqTransaction.Satoshis, events)
	if err != nil {
		h.log.Error().Msgf("An error occurred while creating a transaction: %s", err)
		c.JSON(http.StatusBadRequest, "An error occurred while creating a transaction.")
		return
	}
	go func() {
		transaction := <-events
		h.ws.GetSocket(strconv.Itoa(c.GetInt(auth.SessionUserId))).Notify(transaction)
	}()

	c.Status(http.StatusOK)
}
