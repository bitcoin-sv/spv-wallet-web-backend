package users

import (
	"context"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
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
		SendToRecipients(recipients []*walletclient.Recipients, senderPaymail string) (Transaction, error)
		GetTransactions(queryParam *walletclient.QueryParams, userPaymail string) ([]Transaction, error)
		GetTransaction(transactionId, userPaymail string) (FullTransaction, error)
		GetTransactionsCount() (int64, error)
		CreateAndFinalizeTransaction(recipients []*walletclient.Recipients, metadata *models.Metadata) (DraftTransaction, error)
		RecordTransaction(hex, draftTxId string, metadata *models.Metadata) (*models.Transaction, error)
		// Contacts methods
		UpsertContact(ctx context.Context, paymail, fullName string, metadata *models.Metadata) (*models.Contact, walletclient.ResponseError)
		AcceptContact(ctx context.Context, paymail string) walletclient.ResponseError
		RejectContact(ctx context.Context, paymail string) walletclient.ResponseError
		ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) walletclient.ResponseError
		GetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *walletclient.QueryParams) (*models.SearchContactsResponse, walletclient.ResponseError)
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
