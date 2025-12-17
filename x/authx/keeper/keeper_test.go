package keeper_test

import (
	"testing"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/sagaxyz/ssc/x/authx"
	authxkeeper "github.com/sagaxyz/ssc/x/authx/keeper"
	authxmodule "github.com/sagaxyz/ssc/x/authx/module"
	authxtestutil "github.com/sagaxyz/ssc/x/authx/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	bankSendAuthMsgType = banktypes.SendAuthorization{}.MsgTypeURL()
	coins10             = sdk.NewCoins(sdk.NewInt64Coin("stake", 10))
	coins100            = sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
	coins1000           = sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
)

type TestSuite struct {
	suite.Suite

	ctx           sdk.Context
	addrs         []sdk.AccAddress
	authxKeeper   authxkeeper.Keeper
	accountKeeper *authxtestutil.MockAccountKeeper
	bankKeeper    *authxtestutil.MockBankKeeper
	baseApp       *baseapp.BaseApp
	encCfg        moduletestutil.TestEncodingConfig
	queryClient   authx.QueryClient
	msgSrvr       authx.MsgServer
}

func (s *TestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(authxkeeper.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	s.encCfg = moduletestutil.MakeTestEncodingConfig(authxmodule.AppModuleBasic{})

	s.baseApp = baseapp.NewBaseApp(
		"authx",
		log.NewNopLogger(),
		testCtx.DB,
		s.encCfg.TxConfig.TxDecoder(),
	)
	s.baseApp.SetCMS(testCtx.CMS)
	s.baseApp.SetInterfaceRegistry(s.encCfg.InterfaceRegistry)

	s.addrs = simtestutil.CreateIncrementalAccounts(7)

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	s.accountKeeper = authxtestutil.NewMockAccountKeeper(ctrl)
	s.accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos")).AnyTimes()

	s.bankKeeper = authxtestutil.NewMockBankKeeper(ctrl)
	banktypes.RegisterInterfaces(s.encCfg.InterfaceRegistry)
	banktypes.RegisterMsgServer(s.baseApp.MsgServiceRouter(), s.bankKeeper)
	s.bankKeeper.EXPECT().BlockedAddr(gomock.Any()).Return(false).AnyTimes()

	s.authxKeeper = authxkeeper.NewKeeper(storeService, s.encCfg.Codec, s.baseApp.MsgServiceRouter(), s.accountKeeper).SetBankKeeper(s.bankKeeper)

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.encCfg.InterfaceRegistry)
	authx.RegisterQueryServer(queryHelper, s.authxKeeper)
	queryClient := authx.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.msgSrvr = s.authxKeeper
}

func (s *TestSuite) TestKeeper() {
	ctx, addrs := s.ctx, s.addrs
	now := ctx.BlockTime()
	require := s.Require()

	granterAddr := addrs[0]
	granteeAddr := addrs[1]

	s.T().Log("verify that no authorization returns nil")
	authorizations, err := s.authxKeeper.GetAuthorizations(ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 0)

	s.T().Log("verify save, get and delete")
	sendAutz := &banktypes.SendAuthorization{SpendLimit: coins100}
	expire := now.AddDate(1, 0, 0)
	err = s.authxKeeper.SaveGrant(ctx, granteeAddr, granterAddr, sendAutz, &expire)
	require.NoError(err)

	authorizations, err = s.authxKeeper.GetAuthorizations(ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 1)

	err = s.authxKeeper.DeleteGrant(ctx, granteeAddr, granterAddr, sendAutz.MsgTypeURL())
	require.NoError(err)

	authorizations, err = s.authxKeeper.GetAuthorizations(ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 0)

	s.T().Log("verify granting same authorization overwrite existing authorization")
	err = s.authxKeeper.SaveGrant(ctx, granteeAddr, granterAddr, sendAutz, &expire)
	require.NoError(err)

	authorizations, err = s.authxKeeper.GetAuthorizations(ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 1)

	sendAutz = &banktypes.SendAuthorization{SpendLimit: coins1000}
	err = s.authxKeeper.SaveGrant(ctx, granteeAddr, granterAddr, sendAutz, &expire)
	require.NoError(err)
	authorizations, err = s.authxKeeper.GetAuthorizations(ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 1)
	authorization := authorizations[0]
	sendAuth := authorization.(*banktypes.SendAuthorization)
	require.Equal(sendAuth.SpendLimit, sendAutz.SpendLimit)
	require.Equal(sendAuth.MsgTypeURL(), sendAutz.MsgTypeURL())

	s.T().Log("verify removing non existing authorization returns error")
	err = s.authxKeeper.DeleteGrant(ctx, granterAddr, granteeAddr, "abcd")
	s.Require().Error(err)
}

func (s *TestSuite) TestKeeperIter() {
	ctx, addrs := s.ctx, s.addrs

	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	granter2Addr := addrs[2]
	e := ctx.BlockTime().AddDate(1, 0, 0)
	sendAuthz := banktypes.NewSendAuthorization(coins100, nil)

	s.Require().NoError(s.authxKeeper.SaveGrant(ctx, granteeAddr, granterAddr, sendAuthz, &e))
	s.Require().NoError(s.authxKeeper.SaveGrant(ctx, granteeAddr, granter2Addr, sendAuthz, &e))

	s.authxKeeper.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant authx.Grant) bool {
		s.Require().Equal(granteeAddr, grantee)
		s.Require().Contains([]sdk.AccAddress{granterAddr, granter2Addr}, granter)
		return true
	})
}

func (s *TestSuite) TestDispatchAction() {
	addrs := s.addrs
	require := s.Require()
	now := s.ctx.BlockTime()

	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	recipientAddr := addrs[2]
	a := banktypes.NewSendAuthorization(coins100, nil)

	testCases := []struct {
		name      string
		req       authx.MsgExec
		expectErr bool
		errMsg    string
		preRun    func() sdk.Context
		postRun   func()
	}{
		{
			"expect error authorization not found",
			authx.NewMsgExec(granteeAddr, []sdk.Msg{
				&banktypes.MsgSend{
					Amount:      coins10,
					FromAddress: granterAddr.String(),
					ToAddress:   recipientAddr.String(),
				},
			}),
			true,
			"authorization not found",
			func() sdk.Context {
				// remove any existing authorizations
				_ = s.authxKeeper.DeleteGrant(s.ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
				return s.ctx
			},
			func() {},
		},
		{
			"expect error expired authorization",
			authx.NewMsgExec(granteeAddr, []sdk.Msg{
				&banktypes.MsgSend{
					Amount:      coins10,
					FromAddress: granterAddr.String(),
					ToAddress:   recipientAddr.String(),
				},
			}),
			true,
			"authorization expired",
			func() sdk.Context {
				e := now.AddDate(0, 0, 1)
				err := s.authxKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, a, &e)
				require.NoError(err)
				return s.ctx.WithBlockTime(s.ctx.BlockTime().AddDate(0, 0, 2))
			},
			func() {},
		},
		{
			"expect error over spent limit",
			authx.NewMsgExec(granteeAddr, []sdk.Msg{
				&banktypes.MsgSend{
					Amount:      coins1000,
					FromAddress: granterAddr.String(),
					ToAddress:   recipientAddr.String(),
				},
			}),
			true,
			"requested amount is more than spend limit",
			func() sdk.Context {
				e := now.AddDate(0, 1, 0)
				err := s.authxKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, a, &e)
				require.NoError(err)
				return s.ctx
			},
			func() {},
		},
		{
			"valid test verify amount left",
			authx.NewMsgExec(granteeAddr, []sdk.Msg{
				&banktypes.MsgSend{
					Amount:      coins10,
					FromAddress: granterAddr.String(),
					ToAddress:   recipientAddr.String(),
				},
			}),
			false,
			"",
			func() sdk.Context {
				e := now.AddDate(0, 1, 0)
				err := s.authxKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, a, &e)
				require.NoError(err)
				return s.ctx
			},
			func() {
				authxs, err := s.authxKeeper.GetAuthorizations(s.ctx, granteeAddr, granterAddr)
				require.NoError(err)
				require.Len(authxs, 1)
				authorization := authxs[0].(*banktypes.SendAuthorization)
				require.NotNil(authorization)
				require.Equal(authorization.SpendLimit, coins100.Sub(coins10...))
			},
		},
		{
			"valid test verify authorization is removed when it is used up",
			authx.NewMsgExec(granteeAddr, []sdk.Msg{
				&banktypes.MsgSend{
					Amount:      coins100,
					FromAddress: granterAddr.String(),
					ToAddress:   recipientAddr.String(),
				},
			}),
			false,
			"",
			func() sdk.Context {
				e := now.AddDate(0, 1, 0)
				err := s.authxKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, a, &e)
				require.NoError(err)
				return s.ctx
			},
			func() {
				authxs, err := s.authxKeeper.GetAuthorizations(s.ctx, granteeAddr, granterAddr)
				require.NoError(err)
				require.Len(authxs, 0)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := tc.preRun()
			executeMsgs, err := tc.req.GetMessages()
			require.NoError(err)
			result, err := s.authxKeeper.DispatchActions(ctx, granteeAddr, executeMsgs)
			if tc.expectErr {
				require.Error(err)
				require.Nil(result)
				require.Contains(err.Error(), tc.errMsg)
			} else {
				require.NoError(err)
				require.NotNil(result)
			}
			tc.postRun()
		})
	}
}

// Tests that all msg events included in an authx MsgExec tx
// Ref: https://github.com/cosmos/cosmos-sdk/issues/9501
func (s *TestSuite) TestDispatchedEvents() {
	require := s.Require()
	addrs := s.addrs
	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	recipientAddr := addrs[2]
	expiration := s.ctx.BlockTime().Add(1 * time.Second) // must be in the future

	msgs := authx.NewMsgExec(granteeAddr, []sdk.Msg{
		&banktypes.MsgSend{
			Amount:      coins10,
			FromAddress: granterAddr.String(),
			ToAddress:   recipientAddr.String(),
		},
	})

	// grant authorization
	err := s.authxKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, &banktypes.SendAuthorization{SpendLimit: coins10}, &expiration)
	require.NoError(err)
	authorizations, err := s.authxKeeper.GetAuthorizations(s.ctx, granteeAddr, granterAddr)
	require.NoError(err)
	require.Len(authorizations, 1)
	authorization := authorizations[0].(*banktypes.SendAuthorization)
	require.Equal(authorization.MsgTypeURL(), bankSendAuthMsgType)

	executeMsgs, err := msgs.GetMessages()
	require.NoError(err)

	result, err := s.authxKeeper.DispatchActions(s.ctx, granteeAddr, executeMsgs)
	require.NoError(err)
	require.NotNil(result)

	events := s.ctx.EventManager().Events()

	// get last 5 events (events that occur *after* the grant)
	events = events[len(events)-2:]

	requiredEvents := map[string]bool{
		"cosmos.authx.v1beta1.EventGrant":  true,
		"cosmos.authx.v1beta1.EventRevoke": true,
	}
	for _, e := range events {
		requiredEvents[e.Type] = true
	}
	for _, v := range requiredEvents {
		require.True(v)
	}
}

func (s *TestSuite) TestDequeueAllGrantsQueue() {
	require := s.Require()
	addrs := s.addrs
	granter := addrs[0]
	grantee := addrs[1]
	grantee1 := addrs[2]
	exp := s.ctx.BlockTime().AddDate(0, 0, 1)
	a := banktypes.SendAuthorization{SpendLimit: coins100}

	// create few authorizations
	err := s.authxKeeper.SaveGrant(s.ctx, grantee, granter, &a, &exp)
	require.NoError(err)

	err = s.authxKeeper.SaveGrant(s.ctx, grantee1, granter, &a, &exp)
	require.NoError(err)

	exp2 := exp.AddDate(0, 1, 0)
	err = s.authxKeeper.SaveGrant(s.ctx, granter, grantee1, &a, &exp2)
	require.NoError(err)

	exp2 = exp.AddDate(2, 0, 0)
	err = s.authxKeeper.SaveGrant(s.ctx, granter, grantee, &a, &exp2)
	require.NoError(err)

	newCtx := s.ctx.WithBlockTime(exp.AddDate(1, 0, 0))
	err = s.authxKeeper.DequeueAndDeleteExpiredGrants(newCtx)
	require.NoError(err)

	s.T().Log("verify expired grants are pruned from the state")
	authxs, err := s.authxKeeper.GetAuthorizations(newCtx, grantee, granter)
	require.NoError(err)
	require.Len(authxs, 0)

	authxs, err = s.authxKeeper.GetAuthorizations(newCtx, granter, grantee1)
	require.NoError(err)
	require.Len(authxs, 0)

	authxs, err = s.authxKeeper.GetAuthorizations(newCtx, grantee1, granter)
	require.NoError(err)
	require.Len(authxs, 0)

	authxs, err = s.authxKeeper.GetAuthorizations(newCtx, granter, grantee)
	require.NoError(err)
	require.Len(authxs, 1)
}

func (s *TestSuite) TestGetAuthorization() {
	addr1 := s.addrs[3]
	addr2 := s.addrs[4]
	addr3 := s.addrs[5]
	addr4 := s.addrs[6]

	genAuthMulti := authx.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgMultiSend{}))
	genAuthSend := authx.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{}))
	sendAuth := banktypes.NewSendAuthorization(coins10, nil)

	start := s.ctx.BlockHeader().Time
	expired := start.Add(time.Duration(1) * time.Second)
	notExpired := start.Add(time.Duration(5) * time.Hour)

	s.Require().NoError(s.authxKeeper.SaveGrant(s.ctx, addr1, addr2, genAuthMulti, nil), "creating grant 1->2")
	s.Require().NoError(s.authxKeeper.SaveGrant(s.ctx, addr1, addr3, genAuthSend, &expired), "creating grant 1->3")
	s.Require().NoError(s.authxKeeper.SaveGrant(s.ctx, addr1, addr4, sendAuth, &notExpired), "creating grant 1->4")
	// Without access to private keeper methods, I don't know how to save a grant with an invalid authorization.
	newCtx := s.ctx.WithBlockTime(start.Add(time.Duration(1) * time.Minute))

	tests := []struct {
		name    string
		grantee sdk.AccAddress
		granter sdk.AccAddress
		msgType string
		expAuth authx.Authorization
		expExp  *time.Time
	}{
		{
			name:    "grant has nil exp and is returned",
			grantee: addr1,
			granter: addr2,
			msgType: genAuthMulti.MsgTypeURL(),
			expAuth: genAuthMulti,
			expExp:  nil,
		},
		{
			name:    "grant is expired not returned",
			grantee: addr1,
			granter: addr3,
			msgType: genAuthSend.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
		{
			name:    "grant is not expired and is returned",
			grantee: addr1,
			granter: addr4,
			msgType: sendAuth.MsgTypeURL(),
			expAuth: sendAuth,
			expExp:  &notExpired,
		},
		{
			name:    "grant is not expired but wrong msg type returns nil",
			grantee: addr1,
			granter: addr4,
			msgType: genAuthMulti.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
		{
			name:    "no grant exists between the two",
			grantee: addr2,
			granter: addr3,
			msgType: genAuthSend.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			actAuth, actExp := s.authxKeeper.GetAuthorization(newCtx, tc.grantee, tc.granter, tc.msgType)
			s.Assert().Equal(tc.expAuth, actAuth, "authorization")
			s.Assert().Equal(tc.expExp, actExp, "expiration")
		})
	}
}

func (s *TestSuite) TestGetAuthorizations() {
	require := s.Require()
	addr1 := s.addrs[1]
	addr2 := s.addrs[2]

	genAuthMulti := authx.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgMultiSend{}))
	genAuthSend := authx.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{}))

	start := s.ctx.BlockHeader().Time
	expired := start.Add(time.Duration(1) * time.Second)

	s.Require().NoError(s.authxKeeper.SaveGrant(s.ctx, addr1, addr2, genAuthMulti, &expired), "creating multi send grant 1->2")
	s.Require().NoError(s.authxKeeper.SaveGrant(s.ctx, addr1, addr2, genAuthSend, &expired), "creating send grant 1->2")

	authxs, err := s.authxKeeper.GetAuthorizations(s.ctx, addr1, addr2)
	require.NoError(err)
	require.Len(authxs, 2)
	require.Equal(sdk.MsgTypeURL(&banktypes.MsgMultiSend{}), authxs[0].MsgTypeURL())
	require.Equal(sdk.MsgTypeURL(&banktypes.MsgSend{}), authxs[1].MsgTypeURL())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
