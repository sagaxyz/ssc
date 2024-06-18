package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateChainletStack{}, "chainlet/CreateChainletStack", nil)
	cdc.RegisterConcrete(&MsgLaunchChainlet{}, "chainlet/LaunchChainlet", nil)
	cdc.RegisterConcrete(&MsgUpdateChainletStack{}, "chainlet/UpdateChainletStack", nil)
	cdc.RegisterConcrete(&MsgUpgradeChainlet{}, "chainlet/UpgradeChainlet", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateChainletStack{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgLaunchChainlet{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateChainletStack{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpgradeChainlet{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
