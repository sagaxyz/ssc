package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sagaxyz/ssc/x/escrow/keeper"
	"github.com/sagaxyz/ssc/x/escrow/types"
)

func SimulateMsgWithdraw(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgWithdraw{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the Withdraw simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Withdraw simulation not implemented"), nil, nil
	}
}
