package users

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bip32"
)

type (
	// AccKey is an interface that defianes access key data and methods.
	AccKey interface {
		GetAccessKey() string
		GetAccessKeyId() string
	}
	// PubKey is an interface that defines xpub key data and methods.
	PubKey interface {
		GetId() string
		GetCurrentBalance() uint64
	}

	// Transaction is an interface that defines transaction data and methods.
	Transaction interface {
		GetTransactionId() string
		GetTransactionDirection() string
		GetTransactionTotalValue() uint64
		GetTransactionFee() uint64
		GetTransactionStatus() string
		GetTransactionCreatedDate() time.Time
		GetTransactionSender() string
		GetTransactionReceiver() string
	}

	// FullTransaction is an interface that defines extended transaction data and methods.
	FullTransaction interface {
		GetTransactionId() string
		GetTransactionBlockHash() string
		GetTransactionBlockHeight() uint64
		GetTransactionTotalValue() uint64
		GetTransactionDirection() string
		GetTransactionStatus() string
		GetTransactionFee() uint64
		GetTransactionNumberOfInputs() uint32
		GetTransactionNumberOfOutputs() uint32
		GetTransactionCreatedDate() time.Time
		GetTransactionSender() string
		GetTransactionReceiver() string
	}

	// DraftTransaction is an interface that defines draft transaction data and methods.
	DraftTransaction interface {
		GetDraftTransactionHex() string
		GetDraftTransactionId() string
	}

	// UserWalletClient defines methods which are available for a user with access key.
	UserWalletClient interface {
		// Access Key methods
		CreateAccessKey() (AccKey, error)
		GetAccessKey(accessKeyId string) (AccKey, error)
		RevokeAccessKey(accessKeyId string) (AccKey, error)
		// XPub Key methods
		GetXPub() (PubKey, error)
		// Transaction methods
		SendToRecipients(recipients []*transports.Recipients, senderPaymail string) (Transaction, error)
		GetTransactions(queryParam transports.QueryParams, userPaymail string) ([]Transaction, error)
		GetTransaction(transactionId, userPaymail string) (FullTransaction, error)
		GetTransactionsCount() (int64, error)
		CreateAndFinalizeTransaction(recipients []*transports.Recipients, metadata *models.Metadata) (DraftTransaction, error)
		RecordTransaction(hex, draftTxId string, metadata *models.Metadata) (*models.Transaction, error)
		// Contacts methods
		UpsertContact(ctx context.Context, paymail, fullName string, metadata *models.Metadata) (*models.Contact, transports.ResponseError)
		AcceptContact(ctx context.Context, paymail string) transports.ResponseError
		RejectContact(ctx context.Context, paymail string) transports.ResponseError
		ConfirmContact(ctx context.Context, contact *models.Contact, passcode string, period, digits uint) transports.ResponseError
		GetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *transports.QueryParams) ([]*models.Contact, transports.ResponseError)
		GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error)
	}

	// AdminWalletClient defines methods which are available for an admin with admin key.
	AdminWalletClient interface {
		RegisterXpub(xpriv *bip32.ExtendedKey) (string, error)
		RegisterPaymail(alias, xpub string) (string, error)
		GetSharedConfig() (*models.SharedConfig, error)
	}

	// WalletClientFactory defines methods to create user and admin clients.
	WalletClientFactory interface {
		CreateWithXpriv(xpriv string) (UserWalletClient, error)
		CreateWithAccessKey(accessKey string) (UserWalletClient, error)
		CreateAdminClient() (AdminWalletClient, error)
	}
)
