// Code generated by MockGen. DO NOT EDIT.
// Source: x/billing/types/expected_keepers.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types2 "github.com/cosmos/cosmos-sdk/x/staking/types"
	gomock "github.com/golang/mock/gomock"
	types3 "github.com/sagaxyz/ssc/x/chainlet/types"
	types4 "github.com/sagaxyz/ssc/x/epochs/types"
	types5 "github.com/sagaxyz/ssc/x/escrow/types"
)

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) types0.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(types0.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// GetModuleAccount mocks base method.
func (m *MockAccountKeeper) GetModuleAccount(ctx types.Context, moduleName string) types0.ModuleAccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAccount", ctx, moduleName)
	ret0, _ := ret[0].(types0.ModuleAccountI)
	return ret0
}

// GetModuleAccount indicates an expected call of GetModuleAccount.
func (mr *MockAccountKeeperMockRecorder) GetModuleAccount(ctx, moduleName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAccount), ctx, moduleName)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// GetAccountsBalances mocks base method.
func (m *MockBankKeeper) GetAccountsBalances(ctx types.Context) []types1.Balance {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsBalances", ctx)
	ret0, _ := ret[0].([]types1.Balance)
	return ret0
}

// GetAccountsBalances indicates an expected call of GetAccountsBalances.
func (mr *MockBankKeeperMockRecorder) GetAccountsBalances(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAccountsBalances), ctx)
}

// GetAllBalances mocks base method.
func (m *MockBankKeeper) GetAllBalances(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalances", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// GetAllBalances indicates an expected call of GetAllBalances.
func (mr *MockBankKeeperMockRecorder) GetAllBalances(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAllBalances), ctx, addr)
}

// GetBalance mocks base method.
func (m *MockBankKeeper) GetBalance(ctx types.Context, addr types.AccAddress, denom string) types.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, addr, denom)
	ret0, _ := ret[0].(types.Coin)
	return ret0
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockBankKeeperMockRecorder) GetBalance(ctx, addr, denom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBankKeeper)(nil).GetBalance), ctx, addr, denom)
}

// HasBalance mocks base method.
func (m *MockBankKeeper) HasBalance(ctx types.Context, addr types.AccAddress, amt types.Coin) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasBalance", ctx, addr, amt)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasBalance indicates an expected call of HasBalance.
func (mr *MockBankKeeperMockRecorder) HasBalance(ctx, addr, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasBalance", reflect.TypeOf((*MockBankKeeper)(nil).HasBalance), ctx, addr, amt)
}

// IterateAccountBalances mocks base method.
func (m *MockBankKeeper) IterateAccountBalances(ctx types.Context, addr types.AccAddress, cb func(types.Coin) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateAccountBalances", ctx, addr, cb)
}

// IterateAccountBalances indicates an expected call of IterateAccountBalances.
func (mr *MockBankKeeperMockRecorder) IterateAccountBalances(ctx, addr, cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateAccountBalances", reflect.TypeOf((*MockBankKeeper)(nil).IterateAccountBalances), ctx, addr, cb)
}

// IterateAllBalances mocks base method.
func (m *MockBankKeeper) IterateAllBalances(ctx types.Context, cb func(types.AccAddress, types.Coin) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateAllBalances", ctx, cb)
}

// IterateAllBalances indicates an expected call of IterateAllBalances.
func (mr *MockBankKeeperMockRecorder) IterateAllBalances(ctx, cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).IterateAllBalances), ctx, cb)
}

// LockedCoins mocks base method.
func (m *MockBankKeeper) LockedCoins(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockedCoins", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// LockedCoins indicates an expected call of LockedCoins.
func (mr *MockBankKeeperMockRecorder) LockedCoins(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockedCoins", reflect.TypeOf((*MockBankKeeper)(nil).LockedCoins), ctx, addr)
}

// SendCoinsFromAccountToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx types.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromAccountToModule indicates an expected call of SendCoinsFromAccountToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromAccountToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromAccountToModule), ctx, senderAddr, recipientModule, amt)
}

// SendCoinsFromModuleToAccount mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToAccount", ctx, senderModule, recipientAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToAccount indicates an expected call of SendCoinsFromModuleToAccount.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToAccount", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToAccount), ctx, senderModule, recipientAddr, amt)
}

// SendCoinsFromModuleToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToModule(ctx types.Context, senderModule, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToModule", ctx, senderModule, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToModule indicates an expected call of SendCoinsFromModuleToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToModule), ctx, senderModule, recipientModule, amt)
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), ctx, addr)
}

// ValidateBalance mocks base method.
func (m *MockBankKeeper) ValidateBalance(ctx types.Context, addr types.AccAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBalance", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBalance indicates an expected call of ValidateBalance.
func (mr *MockBankKeeperMockRecorder) ValidateBalance(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBalance", reflect.TypeOf((*MockBankKeeper)(nil).ValidateBalance), ctx, addr)
}

// MockEscrowKeeper is a mock of EscrowKeeper interface.
type MockEscrowKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockEscrowKeeperMockRecorder
}

// MockEscrowKeeperMockRecorder is the mock recorder for MockEscrowKeeper.
type MockEscrowKeeperMockRecorder struct {
	mock *MockEscrowKeeper
}

// NewMockEscrowKeeper creates a new mock instance.
func NewMockEscrowKeeper(ctrl *gomock.Controller) *MockEscrowKeeper {
	mock := &MockEscrowKeeper{ctrl: ctrl}
	mock.recorder = &MockEscrowKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEscrowKeeper) EXPECT() *MockEscrowKeeperMockRecorder {
	return m.recorder
}

// BillAccount mocks base method.
func (m *MockEscrowKeeper) BillAccount(ctx types.Context, amount types.Coin, chainId, toModule string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BillAccount", ctx, amount, chainId, toModule)
	ret0, _ := ret[0].(error)
	return ret0
}

// BillAccount indicates an expected call of BillAccount.
func (mr *MockEscrowKeeperMockRecorder) BillAccount(ctx, amount, chainId, toModule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BillAccount", reflect.TypeOf((*MockEscrowKeeper)(nil).BillAccount), ctx, amount, chainId, toModule)
}

// GetKprChainletAccount mocks base method.
func (m *MockEscrowKeeper) GetKprChainletAccount(ctx types.Context, chainId string) (types5.ChainletAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKprChainletAccount", ctx, chainId)
	ret0, _ := ret[0].(types5.ChainletAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKprChainletAccount indicates an expected call of GetKprChainletAccount.
func (mr *MockEscrowKeeperMockRecorder) GetKprChainletAccount(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKprChainletAccount", reflect.TypeOf((*MockEscrowKeeper)(nil).GetKprChainletAccount), ctx, chainId)
}

// MockStakingKeeper is a mock of StakingKeeper interface.
type MockStakingKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockStakingKeeperMockRecorder
}

// MockStakingKeeperMockRecorder is the mock recorder for MockStakingKeeper.
type MockStakingKeeperMockRecorder struct {
	mock *MockStakingKeeper
}

// NewMockStakingKeeper creates a new mock instance.
func NewMockStakingKeeper(ctrl *gomock.Controller) *MockStakingKeeper {
	mock := &MockStakingKeeper{ctrl: ctrl}
	mock.recorder = &MockStakingKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStakingKeeper) EXPECT() *MockStakingKeeperMockRecorder {
	return m.recorder
}

// GetBondedValidatorsByPower mocks base method.
func (m *MockStakingKeeper) GetBondedValidatorsByPower(ctx types.Context) []types2.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBondedValidatorsByPower", ctx)
	ret0, _ := ret[0].([]types2.Validator)
	return ret0
}

// GetBondedValidatorsByPower indicates an expected call of GetBondedValidatorsByPower.
func (mr *MockStakingKeeperMockRecorder) GetBondedValidatorsByPower(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).GetBondedValidatorsByPower), ctx)
}

// GetValidators mocks base method.
func (m *MockStakingKeeper) GetValidators(ctx types.Context, maxRetrieve uint32) []types2.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, maxRetrieve)
	ret0, _ := ret[0].([]types2.Validator)
	return ret0
}

// GetValidators indicates an expected call of GetValidators.
func (mr *MockStakingKeeperMockRecorder) GetValidators(ctx, maxRetrieve interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockStakingKeeper)(nil).GetValidators), ctx, maxRetrieve)
}

// MockEpochsKeeper is a mock of EpochsKeeper interface.
type MockEpochsKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockEpochsKeeperMockRecorder
}

// MockEpochsKeeperMockRecorder is the mock recorder for MockEpochsKeeper.
type MockEpochsKeeperMockRecorder struct {
	mock *MockEpochsKeeper
}

// NewMockEpochsKeeper creates a new mock instance.
func NewMockEpochsKeeper(ctrl *gomock.Controller) *MockEpochsKeeper {
	mock := &MockEpochsKeeper{ctrl: ctrl}
	mock.recorder = &MockEpochsKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEpochsKeeper) EXPECT() *MockEpochsKeeperMockRecorder {
	return m.recorder
}

// GetEpochInfo mocks base method.
func (m *MockEpochsKeeper) GetEpochInfo(ctx types.Context, identifier string) types4.EpochInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpochInfo", ctx, identifier)
	ret0, _ := ret[0].(types4.EpochInfo)
	return ret0
}

// GetEpochInfo indicates an expected call of GetEpochInfo.
func (mr *MockEpochsKeeperMockRecorder) GetEpochInfo(ctx, identifier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpochInfo", reflect.TypeOf((*MockEpochsKeeper)(nil).GetEpochInfo), ctx, identifier)
}

// MockChainletKeeper is a mock of ChainletKeeper interface.
type MockChainletKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockChainletKeeperMockRecorder
}

// MockChainletKeeperMockRecorder is the mock recorder for MockChainletKeeper.
type MockChainletKeeperMockRecorder struct {
	mock *MockChainletKeeper
}

// NewMockChainletKeeper creates a new mock instance.
func NewMockChainletKeeper(ctrl *gomock.Controller) *MockChainletKeeper {
	mock := &MockChainletKeeper{ctrl: ctrl}
	mock.recorder = &MockChainletKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainletKeeper) EXPECT() *MockChainletKeeperMockRecorder {
	return m.recorder
}

// ChainletExists mocks base method.
func (m *MockChainletKeeper) ChainletExists(ctx types.Context, chainId string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainletExists", ctx, chainId)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ChainletExists indicates an expected call of ChainletExists.
func (mr *MockChainletKeeperMockRecorder) ChainletExists(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainletExists", reflect.TypeOf((*MockChainletKeeper)(nil).ChainletExists), ctx, chainId)
}

// GetChainlet mocks base method.
func (m *MockChainletKeeper) GetChainlet(goCtx context.Context, req *types3.QueryGetChainletRequest) (*types3.QueryGetChainletResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainlet", goCtx, req)
	ret0, _ := ret[0].(*types3.QueryGetChainletResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainlet indicates an expected call of GetChainlet.
func (mr *MockChainletKeeperMockRecorder) GetChainlet(goCtx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainlet", reflect.TypeOf((*MockChainletKeeper)(nil).GetChainlet), goCtx, req)
}

// GetChainletInfo mocks base method.
func (m *MockChainletKeeper) GetChainletInfo(ctx types.Context, chainId string) (*types3.Chainlet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainletInfo", ctx, chainId)
	ret0, _ := ret[0].(*types3.Chainlet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainletInfo indicates an expected call of GetChainletInfo.
func (mr *MockChainletKeeperMockRecorder) GetChainletInfo(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainletInfo", reflect.TypeOf((*MockChainletKeeper)(nil).GetChainletInfo), ctx, chainId)
}

// GetChainletStackInfo mocks base method.
func (m *MockChainletKeeper) GetChainletStackInfo(ctx types.Context, chainId string) (*types3.ChainletStack, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainletStackInfo", ctx, chainId)
	ret0, _ := ret[0].(*types3.ChainletStack)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainletStackInfo indicates an expected call of GetChainletStackInfo.
func (mr *MockChainletKeeperMockRecorder) GetChainletStackInfo(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainletStackInfo", reflect.TypeOf((*MockChainletKeeper)(nil).GetChainletStackInfo), ctx, chainId)
}

// GetParams mocks base method.
func (m *MockChainletKeeper) GetParams(ctx types.Context) types3.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams", ctx)
	ret0, _ := ret[0].(types3.Params)
	return ret0
}

// GetParams indicates an expected call of GetParams.
func (mr *MockChainletKeeperMockRecorder) GetParams(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParams", reflect.TypeOf((*MockChainletKeeper)(nil).GetParams), ctx)
}

// IsChainletStarted mocks base method.
func (m *MockChainletKeeper) IsChainletStarted(ctx types.Context, chainId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChainletStarted", ctx, chainId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsChainletStarted indicates an expected call of IsChainletStarted.
func (mr *MockChainletKeeperMockRecorder) IsChainletStarted(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChainletStarted", reflect.TypeOf((*MockChainletKeeper)(nil).IsChainletStarted), ctx, chainId)
}

// ListChainletStack mocks base method.
func (m *MockChainletKeeper) ListChainletStack(goCtx context.Context, req *types3.QueryListChainletStackRequest) (*types3.QueryListChainletStackResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChainletStack", goCtx, req)
	ret0, _ := ret[0].(*types3.QueryListChainletStackResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChainletStack indicates an expected call of ListChainletStack.
func (mr *MockChainletKeeperMockRecorder) ListChainletStack(goCtx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChainletStack", reflect.TypeOf((*MockChainletKeeper)(nil).ListChainletStack), goCtx, req)
}

// ListChainlets mocks base method.
func (m *MockChainletKeeper) ListChainlets(goCtx context.Context, req *types3.QueryListChainletsRequest) (*types3.QueryListChainletsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChainlets", goCtx, req)
	ret0, _ := ret[0].(*types3.QueryListChainletsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChainlets indicates an expected call of ListChainlets.
func (mr *MockChainletKeeperMockRecorder) ListChainlets(goCtx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChainlets", reflect.TypeOf((*MockChainletKeeper)(nil).ListChainlets), goCtx, req)
}

// StartExistingChainlet mocks base method.
func (m *MockChainletKeeper) StartExistingChainlet(ctx types.Context, chainId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartExistingChainlet", ctx, chainId)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartExistingChainlet indicates an expected call of StartExistingChainlet.
func (mr *MockChainletKeeperMockRecorder) StartExistingChainlet(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartExistingChainlet", reflect.TypeOf((*MockChainletKeeper)(nil).StartExistingChainlet), ctx, chainId)
}

// StopChainlet mocks base method.
func (m *MockChainletKeeper) StopChainlet(ctx types.Context, chainId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopChainlet", ctx, chainId)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopChainlet indicates an expected call of StopChainlet.
func (mr *MockChainletKeeperMockRecorder) StopChainlet(ctx, chainId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopChainlet", reflect.TypeOf((*MockChainletKeeper)(nil).StopChainlet), ctx, chainId)
}