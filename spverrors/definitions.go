package spverrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ////////////////////////////////// CONTACT ERRORS

var ErrUpsertContact = models.SPVError{Message: "Cannot upsert the contact", StatusCode: 400, Code: "error-upsert-contact"}

var ErrAcceptContact = models.SPVError{Message: "Cannot accept the contact", StatusCode: 400, Code: "error-accept-contact"}

var ErrRejectContact = models.SPVError{Message: "Cannot reject the contact", StatusCode: 400, Code: "error-reject-contact"}

var ErrConfirmContact = models.SPVError{Message: "Cannot confirm the contact", StatusCode: 400, Code: "error-confirm-contact"}

var ErrGetContacts = models.SPVError{Message: "Cannot get contacts", StatusCode: 400, Code: "error-get-contacts"}

var ErrGenerateTotpForContact = models.SPVError{Message: "Cannot generate TOTP for contact", StatusCode: 400, Code: "error-generate-totp-for-contact"}

// ////////////////////////////////// TRANSACTION ERRORS

var ErrCreateTransaction = models.SPVError{Message: "Cannot create transaction", StatusCode: 400, Code: "error-transaction-create"}

var ErrGetTransaction = models.SPVError{Message: "Cannot get transaction", StatusCode: 400, Code: "error-transaction-get"}

var ErrGetTransactions = models.SPVError{Message: "Cannot get transactions", StatusCode: 400, Code: "error-transactions-get"}

var ErrCountTransactions = models.SPVError{Message: "Cannot count transactions", StatusCode: 400, Code: "error-transactions-count"}

var ErrRecordTransaction = models.SPVError{Message: "Cannot record transaction", StatusCode: 400, Code: "error-transaction-record"}

// ////////////////////////////////// USER ERRORS

var ErrInvalidCredentials = models.SPVError{Message: "Invalid credentials", StatusCode: 401, Code: "error-credentials-invalid"}

var ErrUserAlreadyExists = models.SPVError{Message: "User already exists", StatusCode: 409, Code: "error-user-already-exists"}

var ErrInsertUser = models.SPVError{Message: "Cannot insert new user", StatusCode: 500, Code: "error-user-insert"}

var ErrEmptyPassword = models.SPVError{Message: "Password cannot be empty", StatusCode: 400, Code: "error-password-empty"}

var ErrIncorrectEmail = models.SPVError{Message: "Incorrect email", StatusCode: 400, Code: "error-email-incorrect"}

var ErrRegisterXPub = models.SPVError{Message: "Cannot register new xPub", StatusCode: 400, Code: "error-xpub-register"}

var ErrRegisterPaymail = models.SPVError{Message: "Cannot register new Paymail", StatusCode: 400, Code: "error-paymail-register"}
