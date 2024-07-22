package chainlet

import (
	"math/rand"

	"github.com/sagaxyz/ssc/testutil/sample"
	chainletsimulation "github.com/sagaxyz/ssc/x/chainlet/simulation"
	"github.com/sagaxyz/ssc/x/chainlet/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = chainletsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateChainletStack = "op_weight_msg_create_chainlet_stack"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateChainletStack int = 100

	opWeightMsgLaunchChainlet = "op_weight_msg_launch_chainlet"
	// TODO: Determine the simulation weight value
	defaultWeightMsgLaunchChainlet int = 100

	opWeightMsgUpdateChainletStack = "op_weight_msg_update_chainlet_stack"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateChainletStack int = 100

	opWeightMsgUpgradeChainlet = "op_weight_msg_upgrade_chainlet"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpgradeChainlet int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	chainletGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&chainletGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {

	return []simtypes.LegacyParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateChainletStack int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateChainletStack, &weightMsgCreateChainletStack, nil,
		func(_ *rand.Rand) {
			weightMsgCreateChainletStack = defaultWeightMsgCreateChainletStack
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateChainletStack,
		chainletsimulation.SimulateMsgCreateChainletStack(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgLaunchChainlet int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgLaunchChainlet, &weightMsgLaunchChainlet, nil,
		func(_ *rand.Rand) {
			weightMsgLaunchChainlet = defaultWeightMsgLaunchChainlet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLaunchChainlet,
		chainletsimulation.SimulateMsgLaunchChainlet(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateChainletStack int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateChainletStack, &weightMsgUpdateChainletStack, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateChainletStack = defaultWeightMsgUpdateChainletStack
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateChainletStack,
		chainletsimulation.SimulateMsgUpdateChainletStack(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpgradeChainlet int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpgradeChainlet, &weightMsgUpgradeChainlet, nil,
		func(_ *rand.Rand) {
			weightMsgUpgradeChainlet = defaultWeightMsgUpgradeChainlet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpgradeChainlet,
		chainletsimulation.SimulateMsgUpgradeChainlet(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
