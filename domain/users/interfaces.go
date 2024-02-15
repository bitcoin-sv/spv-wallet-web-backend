package users

import (
	"time"

	walletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
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

	// UserClient defines methods which are available for a user with access key.
	UserClient interface {
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
		CreateAndFinalizeTransaction(recipients []*transports.Recipients, metadata *walletmodels.Metadata) (DraftTransaction, error)
		RecordTransaction(hex, draftTxId string, metadata *walletmodels.Metadata) (*walletmodels.Transaction, error)
		UnreserveUtxos(draftTxId string) error
	}

	// AdminClient defines methods which are available for an admin with admin key.
	AdminClient interface {
		RegisterXpub(xpriv *bip32.ExtendedKey) (string, error)
		RegisterPaymail(alias, xpub string) (string, error)
	}

	// ClientFactory defines methods to create user and admin clients.
	ClientFactory interface {
		CreateWithXpriv(xpriv string) (UserClient, error)
		CreateWithAccessKey(accessKey string) (UserClient, error)
		CreateAdminClient() (AdminClient, error)
	}
)
