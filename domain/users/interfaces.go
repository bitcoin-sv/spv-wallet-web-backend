package users

import (
	"context"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/libsv/go-bk/bip32"
)

type (
	// AccKey is an interface that defianes access key data and methods.
	AccKey interface {
		GetAccessKey() string
		GetAccessKeyID() string
	}
	// PubKey is an interface that defines xpub key data and methods.
	PubKey interface {
		GetID() string
		GetCurrentBalance() uint64
	}

	// Transaction is an interface that defines transaction data and methods.
	Transaction interface {
		GetTransactionID() string
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
		GetTransactionID() string
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
		GetDraftTransactionID() string
	}

	// UserWalletClient defines methods which are available for a user with access key.
	UserWalletClient interface {
		// Access Key methods
		CreateAccessKey() (AccKey, error)
		GetAccessKey(accessKeyID string) (AccKey, error)
		RevokeAccessKey(accessKeyID string) (AccKey, error)
		// XPub Key methods
		GetXPub() (PubKey, error)
		// Transaction methods
		SendToRecipients(recipients []*walletclient.Recipients, senderPaymail string) (Transaction, error)
		GetTransactions(queryParam *filter.QueryParams, userPaymail string) ([]Transaction, error)
		GetTransaction(transactionID, userPaymail string) (FullTransaction, error)
		GetTransactionsCount() (int64, error)
		CreateAndFinalizeTransaction(recipients []*walletclient.Recipients, metadata map[string]any) (DraftTransaction, error)
		RecordTransaction(hex, draftTxID string, metadata map[string]any) (*models.Transaction, error)
		// Contacts methods
		UpsertContact(ctx context.Context, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error)
		AcceptContact(ctx context.Context, paymail string) error
		RejectContact(ctx context.Context, paymail string) error
		ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error
		GetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error)
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
		CreateWithXpriv(xpriv string) UserWalletClient
		CreateWithAccessKey(accessKey string) UserWalletClient
		CreateAdminClient() (AdminWalletClient, error)
	}
)
