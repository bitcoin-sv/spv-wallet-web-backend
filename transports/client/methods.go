package client

import (
	models "github.com/BuxOrg/bux-models"
)

// GetPaymailsFromMetadata returns sender and receiver paymails from metadata.
// If no paymail was found in metadata, fallback paymail is returned.
func GetPaymailsFromMetadata(transaction *models.Transaction, fallbackPaymail string) (string, string) {
	senderPaymail := ""
	receiverPaymail := ""

	if transaction == nil {
		return senderPaymail, receiverPaymail
	}

	if transaction.Model.Metadata != nil {
		// Try to get paymails from metadata if the transaction was made in SPV Wallet.
		if transaction.Model.Metadata["sender"] != nil {
			senderPaymail = transaction.Model.Metadata["sender"].(string)
		}
		if transaction.Model.Metadata["receiver"] != nil {
			receiverPaymail = transaction.Model.Metadata["receiver"].(string)
		}

		if senderPaymail == "" {
			// Try to get paymails from metadata if the transaction was made outside SPV Wallet.
			if transaction.Model.Metadata["p2p_tx_metadata"] != nil {
				p2pTxMetadata := transaction.Model.Metadata["p2p_tx_metadata"].(map[string]interface{})
				if p2pTxMetadata["sender"] != nil {
					senderPaymail = p2pTxMetadata["sender"].(string)
				}
			}
		}
	}

	if transaction.TransactionDirection == "incoming" && receiverPaymail == "" {
		receiverPaymail = fallbackPaymail
	} else if transaction.TransactionDirection == "outgoing" && senderPaymail == "" {
		senderPaymail = fallbackPaymail
	}

	return senderPaymail, receiverPaymail
}
