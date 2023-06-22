// Code generated by MockGen. DO NOT EDIT.
// Source: domain/users/interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	users "bux-wallet/domain/users"
	reflect "reflect"
	time "time"

	transports "github.com/BuxOrg/go-buxclient/transports"
	gomock "github.com/golang/mock/gomock"
	bip32 "github.com/libsv/go-bk/bip32"
)

// MockAccKey is a mock of AccKey interface.
type MockAccKey struct {
	ctrl     *gomock.Controller
	recorder *MockAccKeyMockRecorder
}

// MockAccKeyMockRecorder is the mock recorder for MockAccKey.
type MockAccKeyMockRecorder struct {
	mock *MockAccKey
}

// NewMockAccKey creates a new mock instance.
func NewMockAccKey(ctrl *gomock.Controller) *MockAccKey {
	mock := &MockAccKey{ctrl: ctrl}
	mock.recorder = &MockAccKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccKey) EXPECT() *MockAccKeyMockRecorder {
	return m.recorder
}

// GetAccessKey mocks base method.
func (m *MockAccKey) GetAccessKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccessKey indicates an expected call of GetAccessKey.
func (mr *MockAccKeyMockRecorder) GetAccessKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessKey", reflect.TypeOf((*MockAccKey)(nil).GetAccessKey))
}

// GetAccessKeyId mocks base method.
func (m *MockAccKey) GetAccessKeyId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessKeyId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccessKeyId indicates an expected call of GetAccessKeyId.
func (mr *MockAccKeyMockRecorder) GetAccessKeyId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessKeyId", reflect.TypeOf((*MockAccKey)(nil).GetAccessKeyId))
}

// MockPubKey is a mock of PubKey interface.
type MockPubKey struct {
	ctrl     *gomock.Controller
	recorder *MockPubKeyMockRecorder
}

// MockPubKeyMockRecorder is the mock recorder for MockPubKey.
type MockPubKeyMockRecorder struct {
	mock *MockPubKey
}

// NewMockPubKey creates a new mock instance.
func NewMockPubKey(ctrl *gomock.Controller) *MockPubKey {
	mock := &MockPubKey{ctrl: ctrl}
	mock.recorder = &MockPubKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubKey) EXPECT() *MockPubKeyMockRecorder {
	return m.recorder
}

// GetCurrentBalance mocks base method.
func (m *MockPubKey) GetCurrentBalance() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentBalance")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetCurrentBalance indicates an expected call of GetCurrentBalance.
func (mr *MockPubKeyMockRecorder) GetCurrentBalance() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentBalance", reflect.TypeOf((*MockPubKey)(nil).GetCurrentBalance))
}

// GetId mocks base method.
func (m *MockPubKey) GetId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetId indicates an expected call of GetId.
func (mr *MockPubKeyMockRecorder) GetId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetId", reflect.TypeOf((*MockPubKey)(nil).GetId))
}

// GetXPub mocks base method.
func (m *MockPubKey) GetXPub() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetXPub")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetXPub indicates an expected call of GetXPub.
func (mr *MockPubKeyMockRecorder) GetXPub() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetXPub", reflect.TypeOf((*MockPubKey)(nil).GetXPub))
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// GetTransactionDirection mocks base method.
func (m *MockTransaction) GetTransactionDirection() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionDirection")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionDirection indicates an expected call of GetTransactionDirection.
func (mr *MockTransactionMockRecorder) GetTransactionDirection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionDirection", reflect.TypeOf((*MockTransaction)(nil).GetTransactionDirection))
}

// GetTransactionId mocks base method.
func (m *MockTransaction) GetTransactionId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionId indicates an expected call of GetTransactionId.
func (mr *MockTransactionMockRecorder) GetTransactionId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionId", reflect.TypeOf((*MockTransaction)(nil).GetTransactionId))
}

// GetTransactionTotalValue mocks base method.
func (m *MockTransaction) GetTransactionTotalValue() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionTotalValue")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTransactionTotalValue indicates an expected call of GetTransactionTotalValue.
func (mr *MockTransactionMockRecorder) GetTransactionTotalValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionTotalValue", reflect.TypeOf((*MockTransaction)(nil).GetTransactionTotalValue))
}

// MockFullTransaction is a mock of FullTransaction interface.
type MockFullTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockFullTransactionMockRecorder
}

// MockFullTransactionMockRecorder is the mock recorder for MockFullTransaction.
type MockFullTransactionMockRecorder struct {
	mock *MockFullTransaction
}

// NewMockFullTransaction creates a new mock instance.
func NewMockFullTransaction(ctrl *gomock.Controller) *MockFullTransaction {
	mock := &MockFullTransaction{ctrl: ctrl}
	mock.recorder = &MockFullTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFullTransaction) EXPECT() *MockFullTransactionMockRecorder {
	return m.recorder
}

// GetTrandsactionCreatedDate mocks base method.
func (m *MockFullTransaction) GetTrandsactionCreatedDate() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrandsactionCreatedDate")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTrandsactionCreatedDate indicates an expected call of GetTrandsactionCreatedDate.
func (mr *MockFullTransactionMockRecorder) GetTrandsactionCreatedDate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrandsactionCreatedDate", reflect.TypeOf((*MockFullTransaction)(nil).GetTrandsactionCreatedDate))
}

// GetTransactionBlockHash mocks base method.
func (m *MockFullTransaction) GetTransactionBlockHash() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionBlockHash")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionBlockHash indicates an expected call of GetTransactionBlockHash.
func (mr *MockFullTransactionMockRecorder) GetTransactionBlockHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionBlockHash", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionBlockHash))
}

// GetTransactionBlockHeight mocks base method.
func (m *MockFullTransaction) GetTransactionBlockHeight() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionBlockHeight")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTransactionBlockHeight indicates an expected call of GetTransactionBlockHeight.
func (mr *MockFullTransactionMockRecorder) GetTransactionBlockHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionBlockHeight", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionBlockHeight))
}

// GetTransactionDirection mocks base method.
func (m *MockFullTransaction) GetTransactionDirection() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionDirection")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionDirection indicates an expected call of GetTransactionDirection.
func (mr *MockFullTransactionMockRecorder) GetTransactionDirection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionDirection", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionDirection))
}

// GetTransactionFee mocks base method.
func (m *MockFullTransaction) GetTransactionFee() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionFee")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTransactionFee indicates an expected call of GetTransactionFee.
func (mr *MockFullTransactionMockRecorder) GetTransactionFee() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionFee", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionFee))
}

// GetTransactionId mocks base method.
func (m *MockFullTransaction) GetTransactionId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionId indicates an expected call of GetTransactionId.
func (mr *MockFullTransactionMockRecorder) GetTransactionId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionId", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionId))
}

// GetTransactionNumberOfInputs mocks base method.
func (m *MockFullTransaction) GetTransactionNumberOfInputs() uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionNumberOfInputs")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// GetTransactionNumberOfInputs indicates an expected call of GetTransactionNumberOfInputs.
func (mr *MockFullTransactionMockRecorder) GetTransactionNumberOfInputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionNumberOfInputs", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionNumberOfInputs))
}

// GetTransactionNumberOfOutputs mocks base method.
func (m *MockFullTransaction) GetTransactionNumberOfOutputs() uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionNumberOfOutputs")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// GetTransactionNumberOfOutputs indicates an expected call of GetTransactionNumberOfOutputs.
func (mr *MockFullTransactionMockRecorder) GetTransactionNumberOfOutputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionNumberOfOutputs", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionNumberOfOutputs))
}

// GetTransactionStatus mocks base method.
func (m *MockFullTransaction) GetTransactionStatus() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionStatus")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTransactionStatus indicates an expected call of GetTransactionStatus.
func (mr *MockFullTransactionMockRecorder) GetTransactionStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionStatus", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionStatus))
}

// GetTransactionTotalValue mocks base method.
func (m *MockFullTransaction) GetTransactionTotalValue() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionTotalValue")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTransactionTotalValue indicates an expected call of GetTransactionTotalValue.
func (mr *MockFullTransactionMockRecorder) GetTransactionTotalValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionTotalValue", reflect.TypeOf((*MockFullTransaction)(nil).GetTransactionTotalValue))
}

// MockUserBuxClient is a mock of UserBuxClient interface.
type MockUserBuxClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserBuxClientMockRecorder
}

// MockUserBuxClientMockRecorder is the mock recorder for MockUserBuxClient.
type MockUserBuxClientMockRecorder struct {
	mock *MockUserBuxClient
}

// NewMockUserBuxClient creates a new mock instance.
func NewMockUserBuxClient(ctrl *gomock.Controller) *MockUserBuxClient {
	mock := &MockUserBuxClient{ctrl: ctrl}
	mock.recorder = &MockUserBuxClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserBuxClient) EXPECT() *MockUserBuxClientMockRecorder {
	return m.recorder
}

// CreateAccessKey mocks base method.
func (m *MockUserBuxClient) CreateAccessKey() (users.AccKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessKey")
	ret0, _ := ret[0].(users.AccKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccessKey indicates an expected call of CreateAccessKey.
func (mr *MockUserBuxClientMockRecorder) CreateAccessKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessKey", reflect.TypeOf((*MockUserBuxClient)(nil).CreateAccessKey))
}

// GetAccessKey mocks base method.
func (m *MockUserBuxClient) GetAccessKey(accessKeyId string) (users.AccKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessKey", accessKeyId)
	ret0, _ := ret[0].(users.AccKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccessKey indicates an expected call of GetAccessKey.
func (mr *MockUserBuxClientMockRecorder) GetAccessKey(accessKeyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessKey", reflect.TypeOf((*MockUserBuxClient)(nil).GetAccessKey), accessKeyId)
}

// GetTransaction mocks base method.
func (m *MockUserBuxClient) GetTransaction(transactionId string) (users.FullTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", transactionId)
	ret0, _ := ret[0].(users.FullTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockUserBuxClientMockRecorder) GetTransaction(transactionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockUserBuxClient)(nil).GetTransaction), transactionId)
}

// GetTransactions mocks base method.
func (m *MockUserBuxClient) GetTransactions() ([]users.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions")
	ret0, _ := ret[0].([]users.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockUserBuxClientMockRecorder) GetTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockUserBuxClient)(nil).GetTransactions))
}

// GetXPub mocks base method.
func (m *MockUserBuxClient) GetXPub() (users.PubKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetXPub")
	ret0, _ := ret[0].(users.PubKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetXPub indicates an expected call of GetXPub.
func (mr *MockUserBuxClientMockRecorder) GetXPub() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetXPub", reflect.TypeOf((*MockUserBuxClient)(nil).GetXPub))
}

// RevokeAccessKey mocks base method.
func (m *MockUserBuxClient) RevokeAccessKey(accessKeyId string) (users.AccKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeAccessKey", accessKeyId)
	ret0, _ := ret[0].(users.AccKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevokeAccessKey indicates an expected call of RevokeAccessKey.
func (mr *MockUserBuxClientMockRecorder) RevokeAccessKey(accessKeyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeAccessKey", reflect.TypeOf((*MockUserBuxClient)(nil).RevokeAccessKey), accessKeyId)
}

// SendToRecipents mocks base method.
func (m *MockUserBuxClient) SendToRecipents(recipients []*transports.Recipients) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendToRecipents", recipients)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendToRecipents indicates an expected call of SendToRecipents.
func (mr *MockUserBuxClientMockRecorder) SendToRecipents(recipients interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendToRecipents", reflect.TypeOf((*MockUserBuxClient)(nil).SendToRecipents), recipients)
}

// MockAdmBuxClient is a mock of AdmBuxClient interface.
type MockAdmBuxClient struct {
	ctrl     *gomock.Controller
	recorder *MockAdmBuxClientMockRecorder
}

// MockAdmBuxClientMockRecorder is the mock recorder for MockAdmBuxClient.
type MockAdmBuxClientMockRecorder struct {
	mock *MockAdmBuxClient
}

// NewMockAdmBuxClient creates a new mock instance.
func NewMockAdmBuxClient(ctrl *gomock.Controller) *MockAdmBuxClient {
	mock := &MockAdmBuxClient{ctrl: ctrl}
	mock.recorder = &MockAdmBuxClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmBuxClient) EXPECT() *MockAdmBuxClientMockRecorder {
	return m.recorder
}

// RegisterPaymail mocks base method.
func (m *MockAdmBuxClient) RegisterPaymail(alias, xpub string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPaymail", alias, xpub)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterPaymail indicates an expected call of RegisterPaymail.
func (mr *MockAdmBuxClientMockRecorder) RegisterPaymail(alias, xpub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPaymail", reflect.TypeOf((*MockAdmBuxClient)(nil).RegisterPaymail), alias, xpub)
}

// RegisterXpub mocks base method.
func (m *MockAdmBuxClient) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterXpub", xpriv)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterXpub indicates an expected call of RegisterXpub.
func (mr *MockAdmBuxClientMockRecorder) RegisterXpub(xpriv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterXpub", reflect.TypeOf((*MockAdmBuxClient)(nil).RegisterXpub), xpriv)
}

// MockBuxClientFactory is a mock of BuxClientFactory interface.
type MockBuxClientFactory struct {
	ctrl     *gomock.Controller
	recorder *MockBuxClientFactoryMockRecorder
}

// MockBuxClientFactoryMockRecorder is the mock recorder for MockBuxClientFactory.
type MockBuxClientFactoryMockRecorder struct {
	mock *MockBuxClientFactory
}

// NewMockBuxClientFactory creates a new mock instance.
func NewMockBuxClientFactory(ctrl *gomock.Controller) *MockBuxClientFactory {
	mock := &MockBuxClientFactory{ctrl: ctrl}
	mock.recorder = &MockBuxClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuxClientFactory) EXPECT() *MockBuxClientFactoryMockRecorder {
	return m.recorder
}

// CreateAdminBuxClient mocks base method.
func (m *MockBuxClientFactory) CreateAdminBuxClient() (users.AdmBuxClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminBuxClient")
	ret0, _ := ret[0].(users.AdmBuxClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdminBuxClient indicates an expected call of CreateAdminBuxClient.
func (mr *MockBuxClientFactoryMockRecorder) CreateAdminBuxClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminBuxClient", reflect.TypeOf((*MockBuxClientFactory)(nil).CreateAdminBuxClient))
}

// CreateWithAccessKey mocks base method.
func (m *MockBuxClientFactory) CreateWithAccessKey(accessKey string) (users.UserBuxClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithAccessKey", accessKey)
	ret0, _ := ret[0].(users.UserBuxClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWithAccessKey indicates an expected call of CreateWithAccessKey.
func (mr *MockBuxClientFactoryMockRecorder) CreateWithAccessKey(accessKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithAccessKey", reflect.TypeOf((*MockBuxClientFactory)(nil).CreateWithAccessKey), accessKey)
}

// CreateWithXpriv mocks base method.
func (m *MockBuxClientFactory) CreateWithXpriv(xpriv string) (users.UserBuxClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithXpriv", xpriv)
	ret0, _ := ret[0].(users.UserBuxClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWithXpriv indicates an expected call of CreateWithXpriv.
func (mr *MockBuxClientFactoryMockRecorder) CreateWithXpriv(xpriv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithXpriv", reflect.TypeOf((*MockBuxClientFactory)(nil).CreateWithXpriv), xpriv)
}
