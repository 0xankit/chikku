package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (

	// TrxsCountKeyPrefix defines the key to store the count of transactions
	TrxsCountKeyPrefix = "TrxsCount/"

	// CumulativeTrxsKeyPrefix defines the key to store the cumulative count of transactions
	CumulativeTrxsKeyPrefix = "CumulativeTrxs/"

	// MemTrxCounterKeyPrefix defines the key to store the count of transactions in memory
	MemTrxCounterKeyPrefix = "MemTrxCounter/"
)

// TrxCountKey returns the store key to retrieve a TrxCounter from the index fields
func TrxCountKey(
	blockHeight int64,
) []byte {
	var key []byte
	blockHeightBytes := sdk.Uint64ToBigEndian(uint64(blockHeight))
	key = append(key, TrxsCountKeyPrefix...)
	key = append(key, blockHeightBytes...)
	key = append(key, []byte("/")...)

	return key
}

// CumulativeTrxsKey returns the store key to retrieve the cumulative TrxCounter
func CumulativeTrxsKey(blockHeight int64) []byte {
	var key []byte
	blockHeightBytes := sdk.Uint64ToBigEndian(uint64(blockHeight))
	key = append(key, CumulativeTrxsKeyPrefix...)
	key = append(key, blockHeightBytes...)
	key = append(key, []byte("/")...)

	return key
}

// MemTrxCounterKey returns the store key to retrieve the MemTrxCounter
func MemTrxCounterKey(blockHeight int64) []byte {
	var key []byte
	blockHeightBytes := sdk.Uint64ToBigEndian(uint64(blockHeight))
	key = append(key, MemTrxCounterKeyPrefix...)
	key = append(key, blockHeightBytes...)
	key = append(key, []byte("/")...)

	return key
}
