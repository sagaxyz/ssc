package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/sagaxyz/ssc/x/authx"
	"github.com/sagaxyz/ssc/x/authx/keeper"
	authxmodule "github.com/sagaxyz/ssc/x/authx/module"
	"github.com/sagaxyz/ssc/x/authx/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestDecodeStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(authxmodule.AppModuleBasic{})
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	dec := simulation.NewDecodeStore(encCfg.Codec)

	now := time.Now().UTC()
	e := now.Add(1)
	sendAuthz := banktypes.NewSendAuthorization(sdk.NewCoins(sdk.NewInt64Coin("foo", 123)), nil)
	grant, _ := authx.NewGrant(now, sendAuthz, &e)
	grantBz, err := encCfg.Codec.Marshal(&grant)
	require.NoError(t, err)
	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: keeper.GrantKey, Value: grantBz},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectErr   bool
		expectedLog string
	}{
		{"Grant", false, fmt.Sprintf("%v\n%v", grant, grant)},
		{"other", true, ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
