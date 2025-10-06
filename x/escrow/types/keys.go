package types

// Module constants
const (
	ModuleName  = "escrow"
	StoreKey    = ModuleName
	RouterKey   = ModuleName
	MemStoreKey = "mem_escrow"
)

// Binary key prefixes (single-byte to keep keys compact).
// Using []byte vars (not const) so we can append at runtime.
var (
	KeyChainletPrefix = []byte{0x01} // escrow/chainlet/{chainId}
	KeyPoolPrefix     = []byte{0x02} // escrow/pool/{chainId}/{denom}
	KeyFunderPrefix   = []byte{0x03} // escrow/funder/{chainId}/{denom}/{addr}
	KeyByFunderPrefix = []byte{0x04} // escrow/byFunder/{addr}/{chainId}/{denom}
)

// delimiter used between path segments. Keep it a single byte.
const delim byte = '/'

// --- Chainlet keys ---

// ChainletKey -> escrow/chainlet/{chainId}
func ChainletKey(chainID string) []byte {
	k := make([]byte, 0, 1+len(chainID))
	k = append(k, KeyChainletPrefix...)
	k = append(k, []byte(chainID)...)
	return k
}

// --- Pool keys (per {chainId, denom}) ---

// PoolKey -> escrow/pool/{chainId}/{denom}
func PoolKey(chainID, denom string) []byte {
	k := make([]byte, 0, 2+len(chainID)+len(denom))
	k = append(k, KeyPoolPrefix...)
	k = append(k, []byte(chainID)...)
	k = append(k, delim)
	k = append(k, []byte(denom)...)
	return k
}

// PoolPrefix -> escrow/pool/{chainId}/  (for iterating all pools of a chainlet)
func PoolPrefix(chainID string) []byte {
	k := make([]byte, 0, 2+len(chainID))
	k = append(k, KeyPoolPrefix...)
	k = append(k, []byte(chainID)...)
	k = append(k, delim)
	return k
}

// --- Funder keys (per {chainId, denom, addr}) ---

// FunderKey -> escrow/funder/{chainId}/{denom}/{addr}
func FunderKey(chainID, denom, addr string) []byte {
	k := FunderPrefix(chainID, denom)
	k = append(k, []byte(addr)...)
	return k
}

// FunderPrefix -> escrow/funder/{chainId}/{denom}/  (for iterating funders in a pool)
func FunderPrefix(chainID, denom string) []byte {
	k := make([]byte, 0, 3+len(chainID)+len(denom))
	k = append(k, KeyFunderPrefix...)
	k = append(k, []byte(chainID)...)
	k = append(k, delim)
	k = append(k, []byte(denom)...)
	k = append(k, delim)
	return k
}

// --- Reverse index (for "list my positions") ---

// ByFunderKey -> escrow/byFunder/{addr}/{chainId}/{denom}
func ByFunderKey(addr, chainID, denom string) []byte {
	k := make([]byte, 0, 3+len(addr)+len(chainID)+len(denom))
	k = append(k, KeyByFunderPrefix...)
	k = append(k, []byte(addr)...)
	k = append(k, delim)
	k = append(k, []byte(chainID)...)
	k = append(k, delim)
	k = append(k, []byte(denom)...)
	return k
}

// ByFunderPrefix -> escrow/byFunder/{addr}/  (for iterating a wallet's positions)
func ByFunderPrefix(addr string) []byte {
	k := make([]byte, 0, 2+len(addr))
	k = append(k, KeyByFunderPrefix...)
	k = append(k, []byte(addr)...)
	k = append(k, delim)
	return k
}
