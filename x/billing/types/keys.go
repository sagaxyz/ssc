package types

const (
	// ModuleName defines the module name
	ModuleName = "billing"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_billing"

	AttributeKeyAccount = "account"
	AttributeKeyAmount  = "amount"
	AttributeKeyMemo    = "memo"
)

var (
	BillingHistoryKey         = "billinghistory"
	ValidatorPayoutHistoryKey = "validatorpayout"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
