package main

import (
	"os"

	//"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sagaxyz/ssc/app"
	"github.com/sagaxyz/ssc/cmd/sscd/cmd"
)

func main() {
	cfg := sdk.GetConfig()
	cmd.SetBech32Prefixes(cfg)

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		//case server.ErrorCode:
		//	os.Exit(e.Code) //TODO
		default:
			_ = e
			os.Exit(1)
		}
	}
}
