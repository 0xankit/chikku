package types

const (
	// ModuleName defines the module name
	ModuleName = "egvmod"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_egvmod"
)

var (
	ParamsKey = []byte("p_egvmod")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
