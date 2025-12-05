package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"

	cosmossdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sagaxyz/ssc/x/peers/types"
)

var hostAllowed = regexp.MustCompile(`^[A-Za-z0-9\.\-\[\]:]+$`) // hostname or IP

// Basic validation of the ID@addr:port format.
// Only needs the bare minimum check for safety.
func validateAddress(addr string) error {
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return errors.New("missing @")
	}
	id := parts[0]
	hostPort := parts[1]

	// Check ID is a hex string
	if len(id) == 0 {
		return errors.New("missing ID")
	}
	_, err := hex.DecodeString(id)
	if err != nil {
		return err
	}

	// Check valid host:port
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(port) // port is a number
	if err != nil {
		return err
	}
	if !hostAllowed.MatchString(host) { // only allowed characters
		return errors.New("invalid characters in host")
	}

	return nil
}

func (k msgServer) SetPeers(goCtx context.Context, msg *types.MsgSetPeers) (resp *types.MsgSetPeersResponse, err error) {
	err = msg.ValidateBasic()
	if err != nil {
		return
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err = k.chainletKeeper.Chainlet(ctx, msg.ChainId)
	if err != nil {
		err = errors.New("no such chain ID")
		return
	}

	if len(msg.Peers) == 0 {
		err = errors.New("no peers provided")
		return
	}
	p := k.GetParams(ctx)
	var dataSize uint32
	for _, addr := range msg.Peers {
		if len(addr) > math.MaxUint32 {
			err = errors.New("data size exceeds uint32")
			return
		}
		dataSize += uint32(len(addr))
		if dataSize > p.MaxData {
			err = fmt.Errorf("exceeded maximum size (%d) of peers", p.MaxData)
			return
		}
		err = validateAddress(addr)
		if err != nil {
			err = fmt.Errorf("invalid addr '%s' in peers: %w", addr, err)
			return
		}
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		err = cosmossdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
		return
	}
	valAddr := sdk.ValAddress(accAddr)

	data := types.Data{
		Updated:   ctx.BlockTime(),
		Addresses: msg.Peers,
	}
	k.StoreData(ctx, msg.ChainId, valAddr.String(), data)

	err = ctx.EventManager().EmitTypedEvent(&types.EventUpdatedChainlet{
		ChainId: msg.ChainId,
	})
	if err != nil {
		return
	}

	resp = &types.MsgSetPeersResponse{}
	return
}
