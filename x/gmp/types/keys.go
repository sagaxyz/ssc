package types

import "strconv"

const (
	// ModuleName defines the module name
	ModuleName = "gmp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_gmp"

	// Version defines the current version the IBC module supports
	Version = "gmp-1"

	// PortID is the default port id that module binds to
	PortID = "gmp"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey      = KeyPrefix("gmp-port-")
	PacketPrefix = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PacketKey(channelID, port string, sequence uint64) []byte {
	k := make([]byte, 0, len(PacketPrefix)+len(channelID)+len(port)+8)
	k = append(k, PacketPrefix...)
	k = append(k, []byte(channelID)...)
	k = append(k, []byte(port)...)
	k = append(k, []byte(strconv.FormatUint(sequence, 10))...)
	return k
}
