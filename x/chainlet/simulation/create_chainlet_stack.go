package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sagaxyz/ssc/x/chainlet/keeper"
	"github.com/sagaxyz/ssc/x/chainlet/types"
)

func SimulateMsgCreateChainletStack(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateChainletStack{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the CreateChainletStack simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "CreateChainletStack simulation not implemented"), nil, nil
	}
}
