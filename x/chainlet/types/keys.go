package types

const (
	// ModuleName defines the module name
	ModuleName = "chainlet"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_chainlet"

	// Version defines the current version the IBC module supports
	Version = "chainlet-1"

	// PortID is the default port id that module binds to
	PortID = "chainlet"
)

var (
	ChainletKey           = []byte{0x01}
	ChainletStackKey      = []byte{0x02}
	ChainletInit          = []byte{0x03}
	NumChainletsKey       = []byte{0x04}
	UpgradingChainletsKey = []byte{0x05}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("chainlet-port-")
)
