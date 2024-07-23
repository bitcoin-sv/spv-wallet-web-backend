package spverrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ////////////////////////////////// AUTHORIZATION ERRORS

var ErrUpsertContact = models.SPVError{Message: "Cannot upsert the contact", StatusCode: 400, Code: "error-upsert-contact"}

var ErrAcceptContact = models.SPVError{Message: "Cannot accept the contact", StatusCode: 400, Code: "error-accept-contact"}

var ErrRejectContact = models.SPVError{Message: "Cannot reject the contact", StatusCode: 400, Code: "error-reject-contact"}

var ErrConfirmContact = models.SPVError{Message: "Cannot confirm the contact", StatusCode: 400, Code: "error-confirm-contact"}

var ErrGetContacts = models.SPVError{Message: "Cannot get contacts", StatusCode: 400, Code: "error-get-contacts"}

var ErrGenerateTotpForContact = models.SPVError{Message: "Cannot generate TOTP for contact", StatusCode: 400, Code: "error-generate-totp-for-contact"}
