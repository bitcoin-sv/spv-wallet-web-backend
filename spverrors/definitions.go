package spverrors

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// ////////////////////////////////// CONTACT ERRORS

// ErrUpsertContact indicates failure to upsert the contact
var ErrUpsertContact = models.SPVError{
	Message:    "Cannot upsert the contact",
	StatusCode: http.StatusBadRequest,
	Code:       "error-upsert-contact",
}

// ErrAcceptContact indicates failure to accept the contact
var ErrAcceptContact = models.SPVError{
	Message:    "Cannot accept the contact",
	StatusCode: http.StatusBadRequest,
	Code:       "error-accept-contact",
}

// ErrRejectContact indicates failure to reject the contact
var ErrRejectContact = models.SPVError{
	Message:    "Cannot reject the contact",
	StatusCode: http.StatusBadRequest,
	Code:       "error-reject-contact",
}

// ErrConfirmContact indicates failure to confirm the contact
var ErrConfirmContact = models.SPVError{
	Message:    "Cannot confirm the contact",
	StatusCode: http.StatusBadRequest,
	Code:       "error-confirm-contact",
}

// ErrGetContacts indicates failure to get contacts
var ErrGetContacts = models.SPVError{
	Message:    "Cannot get contacts",
	StatusCode: http.StatusBadRequest,
	Code:       "error-get-contacts",
}

// ErrGenerateTotpForContact indicates failure to generate TOTP for contact
var ErrGenerateTotpForContact = models.SPVError{
	Message:    "Cannot generate TOTP for contact",
	StatusCode: http.StatusBadRequest,
	Code:       "error-generate-totp-for-contact",
}

// ErrContactNotProvided indicates the contact was not provided
var ErrContactNotProvided = models.SPVError{
	Message:    "Contact not provided",
	StatusCode: http.StatusBadRequest,
	Code:       "error-contact-not-provided",
}

// ////////////////////////////////// TRANSACTION ERRORS

// ErrCreateTransaction indicates failure to create a transaction
var ErrCreateTransaction = models.SPVError{
	Message:    "Cannot create transaction",
	StatusCode: http.StatusBadRequest,
	Code:       "error-transaction-create",
}

// ErrGetTransaction indicates failure to get a transaction
var ErrGetTransaction = models.SPVError{
	Message:    "Cannot get transaction",
	StatusCode: http.StatusBadRequest,
	Code:       "error-transaction-get",
}

// ErrGetTransactions indicates failure to get transactions
var ErrGetTransactions = models.SPVError{
	Message:    "Cannot get transactions",
	StatusCode: http.StatusBadRequest,
	Code:       "error-transactions-get",
}

// ErrCountTransactions indicates failure to count transactions
var ErrCountTransactions = models.SPVError{
	Message:    "Cannot count transactions",
	StatusCode: http.StatusBadRequest,
	Code:       "error-transactions-count",
}

// ErrRecordTransaction indicates failure to record a transaction
var ErrRecordTransaction = models.SPVError{
	Message:    "Cannot record transaction",
	StatusCode: http.StatusBadRequest,
	Code:       "error-transaction-record",
}

// ////////////////////////////////// USER ERRORS

// ErrUnauthorized indicates the user is unauthorized
var ErrUnauthorized = models.SPVError{
	Message:    "Unauthorized",
	StatusCode: http.StatusUnauthorized,
	Code:       "error-unauthorized",
}

// ErrInvalidCredentials indicates invalid credentials were provided
var ErrInvalidCredentials = models.SPVError{
	Message:    "Invalid credentials",
	StatusCode: http.StatusUnauthorized,
	Code:       "error-credentials-invalid",
}

// ErrUserAlreadyExists indicates the user already exists
var ErrUserAlreadyExists = models.SPVError{
	Message:    "User already exists",
	StatusCode: http.StatusConflict,
	Code:       "error-user-already-exists",
}

// ErrInsertUser indicates failure to insert a new user
var ErrInsertUser = models.SPVError{
	Message:    "Cannot insert new user",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-user-insert",
}

// ErrEmptyPassword indicates the password cannot be empty
var ErrEmptyPassword = models.SPVError{
	Message:    "Password cannot be empty",
	StatusCode: http.StatusBadRequest,
	Code:       "error-password-empty",
}

// ErrPasswordMismatch indicates the password and confirmation password do not match
var ErrPasswordMismatch = models.SPVError{
	Message:    "Password and confirmation password do not match",
	StatusCode: http.StatusBadRequest,
	Code:       "error-password-mismatch",
}

// ErrIncorrectEmail indicates an incorrect email was provided
var ErrIncorrectEmail = models.SPVError{
	Message:    "Incorrect email",
	StatusCode: http.StatusBadRequest,
	Code:       "error-email-incorrect",
}

// ErrRegisterXPub indicates failure to register a new xPub
var ErrRegisterXPub = models.SPVError{
	Message:    "Cannot register new xPub",
	StatusCode: http.StatusBadRequest,
	Code:       "error-xpub-register",
}

// ErrRegisterPaymail indicates failure to register a new Paymail
var ErrRegisterPaymail = models.SPVError{
	Message:    "Cannot register new Paymail",
	StatusCode: http.StatusBadRequest,
	Code:       "error-paymail-register",
}

// ErrGenerateMnemonic indicates failure to generate a mnemonic
var ErrGenerateMnemonic = models.SPVError{
	Message:    "Cannot generate mnemonic",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-mnemonic-generate",
}

// ErrGenerateXPriv indicates failure to generate an xPriv
var ErrGenerateXPriv = models.SPVError{
	Message:    "Cannot generate xPriv",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-xpriv-generate",
}

// ErrEncryptXPriv indicates failure to encrypt an xPriv
var ErrEncryptXPriv = models.SPVError{
	Message:    "Cannot encrypt xPriv",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-xpriv-encrypt",
}

// ErrGetUser indicates failure to get user information
var ErrGetUser = models.SPVError{
	Message:    "Cannot get user",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-user-get",
}

// ErrCreateAccessKey indicates failure to create an access key
var ErrCreateAccessKey = models.SPVError{
	Message:    "Cannot create access key",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-access-key-create",
}

// ErrGetXPub indicates failure to get an xPub
var ErrGetXPub = models.SPVError{
	Message:    "Cannot get xPub",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-xpub-get",
}

// ErrGetBalance indicates failure to get the balance
var ErrGetBalance = models.SPVError{
	Message:    "Cannot get balance",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-balance-get",
}

// ErrSessionUpdate indicates failure to update the session
var ErrSessionUpdate = models.SPVError{
	Message:    "Cannot update session",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-session-update",
}

// ErrSessionTerminate indicates failure to terminate the session
var ErrSessionTerminate = models.SPVError{
	Message:    "Cannot terminate session",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-session-terminate",
}

// ////////////////////////////////// RATE ERRORS

// ErrRateNotFound indicates the requested rate was not found
var ErrRateNotFound = models.SPVError{
	Message:    "Rate not found",
	StatusCode: http.StatusNotFound,
	Code:       "error-rate-not-found",
}

// ////////////////////////////////// BINDING ERRORS

// ErrCannotBindRequest is when request body cannot be bind into struct
var ErrCannotBindRequest = models.SPVError{
	Message:    "cannot bind request body",
	StatusCode: 400,
	Code:       "error-bind-body-invalid",
}

// ////////////////////////////////// CONFIG ERRORS

// ErrGetConfig indicates failure to get the configuration
var ErrGetConfig = models.SPVError{
	Message:    "Cannot get configuration",
	StatusCode: http.StatusInternalServerError,
	Code:       "error-config-get",
}
