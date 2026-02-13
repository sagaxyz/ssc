package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	evmhd "github.com/cosmos/evm/crypto/hd"
)

func WithEthereum(options *keyring.Options) {
	options.SupportedAlgos = append(options.SupportedAlgos, evmhd.EthSecp256k1)
	options.SupportedAlgosLedger = append(options.SupportedAlgosLedger, evmhd.EthSecp256k1)
}
