package keeper_test

import (
	"time"

	"go.uber.org/mock/gomock"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec/address"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *TestSuite) createAccounts(accs int) []sdk.AccAddress {
	addrs := simtestutil.CreateIncrementalAccounts(2)
	suite.accountKeeper.EXPECT().GetAccount(gomock.Any(), suite.addrs[0]).Return(authtypes.NewBaseAccountWithAddress(suite.addrs[0])).AnyTimes()
	suite.accountKeeper.EXPECT().GetAccount(gomock.Any(), suite.addrs[1]).Return(authtypes.NewBaseAccountWithAddress(suite.addrs[1])).AnyTimes()
	return addrs
}

func (suite *TestSuite) TestGrant() {
	ctx := suite.ctx.WithBlockTime(time.Now())
	addrs := suite.createAccounts(2)
	curBlockTime := ctx.BlockTime()

	suite.accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos")).AnyTimes()

	oneHour := curBlockTime.Add(time.Hour)
	oneYear := curBlockTime.AddDate(1, 0, 0)

	coins := sdk.NewCoins(sdk.NewCoin("steak", sdkmath.NewInt(10)))

	grantee, granter := addrs[0], addrs[1]

	testCases := []struct {
		name     string
		malleate func() *authx.MsgGrant
		expErr   bool
		errMsg   string
	}{
		{
			name: "identical grantee and granter",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: grantee.String(),
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
			expErr: true,
			errMsg: "grantee and granter should be different",
		},
		{
			name: "invalid granter",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: "invalid",
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
			expErr: true,
			errMsg: "invalid bech32 string",
		},
		{
			name: "invalid grantee",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: "invalid",
					Grant:   grant,
				}
			},
			expErr: true,
			errMsg: "invalid bech32 string",
		},
		{
			name: "invalid grant",
			malleate: func() *authx.MsgGrant {
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant: authx.Grant{
						Expiration: &oneYear,
					},
				}
			},
			expErr: true,
			errMsg: "authorization is nil: invalid type",
		},
		{
			name: "invalid grant, past time",
			malleate: func() *authx.MsgGrant {
				pastTime := curBlockTime.Add(-time.Hour)
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneHour) // we only need the authorization
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant: authx.Grant{
						Authorization: grant.Authorization,
						Expiration:    &pastTime,
					},
				}
			},
			expErr: true,
			errMsg: "expiration must be after the current block time",
		},
		{
			name: "grantee account does not exist on chain: valid grant",
			malleate: func() *authx.MsgGrant {
				newAcc := sdk.AccAddress("valid")
				suite.accountKeeper.EXPECT().GetAccount(gomock.Any(), newAcc).Return(nil).AnyTimes()
				acc := authtypes.NewBaseAccountWithAddress(newAcc)
				suite.accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), newAcc).Return(acc).AnyTimes()
				suite.accountKeeper.EXPECT().SetAccount(gomock.Any(), acc).Return()

				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: newAcc.String(),
					Grant:   grant,
				}
			},
		},
		{
			name: "valid grant",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
		},
		{
			name: "valid grant, same grantee, granter pair but different msgType",
			malleate: func() *authx.MsgGrant {
				g, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), &oneHour)
				suite.Require().NoError(err)
				_, err = suite.msgSrvr.Grant(suite.ctx, &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant:   g,
				})
				suite.Require().NoError(err)

				grant, err := authx.NewGrant(curBlockTime, authx.NewGenericAuthorization("/cosmos.bank.v1beta1.MsgUpdateParams"), &oneHour)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
		},
		{
			name: "valid grant with allow list",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, []sdk.AccAddress{granter}), &oneYear)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
		},
		{
			name: "valid grant with nil expiration time",
			malleate: func() *authx.MsgGrant {
				grant, err := authx.NewGrant(curBlockTime, banktypes.NewSendAuthorization(coins, nil), nil)
				suite.Require().NoError(err)
				return &authx.MsgGrant{
					Granter: granter.String(),
					Grantee: grantee.String(),
					Grant:   grant,
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgSrvr.Grant(suite.ctx, tc.malleate())
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *TestSuite) TestRevoke() {
	addrs := suite.createAccounts(2)

	grantee, granter := addrs[0], addrs[1]

	testCases := []struct {
		name     string
		malleate func() *authx.MsgRevoke
		expErr   bool
		errMsg   string
	}{
		{
			name: "identical grantee and granter",
			malleate: func() *authx.MsgRevoke {
				return &authx.MsgRevoke{
					Granter:    grantee.String(),
					Grantee:    grantee.String(),
					MsgTypeUrl: bankSendAuthMsgType,
				}
			},
			expErr: true,
			errMsg: "grantee and granter should be different",
		},
		{
			name: "invalid granter",
			malleate: func() *authx.MsgRevoke {
				return &authx.MsgRevoke{
					Granter:    "invalid",
					Grantee:    grantee.String(),
					MsgTypeUrl: bankSendAuthMsgType,
				}
			},
			expErr: true,
			errMsg: "invalid bech32 string",
		},
		{
			name: "invalid grantee",
			malleate: func() *authx.MsgRevoke {
				return &authx.MsgRevoke{
					Granter:    granter.String(),
					Grantee:    "invalid",
					MsgTypeUrl: bankSendAuthMsgType,
				}
			},
			expErr: true,
			errMsg: "invalid bech32 string",
		},
		{
			name: "no msg given",
			malleate: func() *authx.MsgRevoke {
				return &authx.MsgRevoke{
					Granter:    granter.String(),
					Grantee:    grantee.String(),
					MsgTypeUrl: "",
				}
			},
			expErr: true,
			errMsg: "missing msg method name",
		},
		{
			name: "valid grant",
			malleate: func() *authx.MsgRevoke {
				suite.createSendAuthorization(grantee, granter)

				return &authx.MsgRevoke{
					Granter:    granter.String(),
					Grantee:    grantee.String(),
					MsgTypeUrl: bankSendAuthMsgType,
				}
			},
		},
		{
			name: "no existing grant to revoke",
			malleate: func() *authx.MsgRevoke {
				return &authx.MsgRevoke{
					Granter:    granter.String(),
					Grantee:    grantee.String(),
					MsgTypeUrl: bankSendAuthMsgType,
				}
			},
			expErr: true,
			errMsg: "authorization not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgSrvr.Revoke(suite.ctx, tc.malleate())
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *TestSuite) TestExec() {
	addrs := suite.createAccounts(2)

	grantee, granter := addrs[0], addrs[1]
	coins := sdk.NewCoins(sdk.NewCoin("steak", sdkmath.NewInt(10)))

	msg := &banktypes.MsgSend{
		FromAddress: granter.String(),
		ToAddress:   grantee.String(),
		Amount:      coins,
	}

	testCases := []struct {
		name     string
		malleate func() authx.MsgExec
		expErr   bool
		errMsg   string
	}{
		{
			name: "invalid grantee (empty)",
			malleate: func() authx.MsgExec {
				return authx.NewMsgExec(sdk.AccAddress{}, []sdk.Msg{msg})
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "non existing grant",
			malleate: func() authx.MsgExec {
				return authx.NewMsgExec(grantee, []sdk.Msg{msg})
			},
			expErr: true,
			errMsg: "authorization not found",
		},
		{
			name: "no message case",
			malleate: func() authx.MsgExec {
				return authx.NewMsgExec(grantee, []sdk.Msg{})
			},
			expErr: true,
			errMsg: "messages cannot be empty",
		},
		{
			name: "valid case",
			malleate: func() authx.MsgExec {
				suite.createSendAuthorization(grantee, granter)
				return authx.NewMsgExec(grantee, []sdk.Msg{msg})
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req := tc.malleate()
			_, err := suite.msgSrvr.Exec(suite.ctx, &req)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
