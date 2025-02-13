package keeper

import (
	"chikku/x/egvmod/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTrxCount sets the transaction count at the given block height
func (k Keeper) SetTrxCount(ctx sdk.Context, operatorsTrxsCount types.OperatorsTrxsCount) {
	// validate operatorsTrxsCount
	if operatorsTrxsCount.BlockHeight != ctx.BlockHeight() {
		return
	}
	// set total transaction count
	var totalTrxCount int64
	for _, operatorTrxCount := range operatorsTrxsCount.OperatorTrxCounters {
		totalTrxCount += operatorTrxCount.TrxCount
	}
	operatorsTrxsCount.TrxCount = totalTrxCount
	// store the transaction count
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TrxsCountKeyPrefix))
	b := k.cdc.MustMarshal(&operatorsTrxsCount)

	store.Set(types.TrxCountKey(
		operatorsTrxsCount.BlockHeight,
	), b)

	// set cumulative transaction count
	k.setCumulativeTrxCount(ctx, operatorsTrxsCount)
}

// setCumulativeTrxCounter sets the cumulative transaction count
func (k Keeper) setCumulativeTrxCount(ctx sdk.Context, operatorsTrxsCount types.OperatorsTrxsCount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CumulativeTrxsKeyPrefix))

	// Retrieve the previous cumulative transaction count (from BlockHeight - 1)
	prevCumulativeTrxCount := k.GetCumulativeTrxs(ctx, operatorsTrxsCount.BlockHeight-1)

	// Create a new cumulative transaction count
	var cumulativeTrxCount types.OperatorsTrxsCount
	cumulativeTrxCount.BlockHeight = operatorsTrxsCount.BlockHeight
	cumulativeTrxCount.TrxCount = prevCumulativeTrxCount.TrxCount + operatorsTrxsCount.TrxCount

	// Create a map to track cumulative transactions per operator
	operatorTrxMap := make(map[string]int64)

	// Add previous cumulative operator transaction counts
	for _, prevOperatorTrx := range prevCumulativeTrxCount.OperatorTrxCounters {
		operatorTrxMap[prevOperatorTrx.Operator] = prevOperatorTrx.TrxCount
	}

	// Update with new transaction counts from the current block
	for _, operatorTrxCount := range operatorsTrxsCount.OperatorTrxCounters {
		operatorTrxMap[operatorTrxCount.Operator] += operatorTrxCount.TrxCount
	}

	// Convert map back to `OperatorTrxCounter` list
	for operator, trxCount := range operatorTrxMap {
		cumulativeTrxCount.OperatorTrxCounters = append(cumulativeTrxCount.OperatorTrxCounters, &types.OperatorTrxCounter{
			Operator: operator,
			TrxCount: trxCount,
		})
	}

	// Marshal and store the updated cumulative transaction count
	b := k.cdc.MustMarshal(&cumulativeTrxCount)
	store.Set(types.CumulativeTrxsKey(cumulativeTrxCount.BlockHeight), b)
}

// GetCumulativeTrxs returns the current cumulative transaction count
func (k Keeper) GetCumulativeTrxs(ctx sdk.Context, blockHeight int64) (cumulativeTrxCount types.OperatorsTrxsCount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CumulativeTrxsKeyPrefix))
	b := store.Get(types.CumulativeTrxsKey(blockHeight))
	if b == nil {
		return types.OperatorsTrxsCount{}
	}
	k.cdc.MustUnmarshal(b, &cumulativeTrxCount)
	return cumulativeTrxCount
}

// GetTrxCount returns the transaction count at the given block height
func (k Keeper) GetTrxCount(ctx sdk.Context, blockHeight int64) (operatorsTrxsCount types.OperatorsTrxsCount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TrxsCountKeyPrefix))
	b := store.Get(types.TrxCountKey(blockHeight))
	if b == nil {
		return types.OperatorsTrxsCount{}
	}
	k.cdc.MustUnmarshal(b, &operatorsTrxsCount)
	return operatorsTrxsCount
}

// GetTrxCountsInRange returns transaction counts for each operator between two block heights
func (k Keeper) GetTrxCountsInRange(ctx sdk.Context, lastDistributionHeight, distributionHeight int64) (operatorsTrxsCounts []types.OperatorTrxCounter, totalTrxCount int64) {
	if lastDistributionHeight > distributionHeight {
		return operatorsTrxsCounts, totalTrxCount
	}

	// Retrieve cumulative transactions at both block heights
	endBlockCumulativeTrxs := k.GetCumulativeTrxs(ctx, distributionHeight)
	startBlockCumulativeTrxs := k.GetCumulativeTrxs(ctx, lastDistributionHeight-1)

	// Create a map for quick lookup of starting block transaction counts
	startTrxMap := make(map[string]int64)
	for _, operatorTrxCounter := range startBlockCumulativeTrxs.OperatorTrxCounters {
		startTrxMap[operatorTrxCounter.Operator] = operatorTrxCounter.TrxCount
	}

	// Compute the difference in transaction counts
	for _, operatorTrxCounter := range endBlockCumulativeTrxs.OperatorTrxCounters {
		startTrxCount := startTrxMap[operatorTrxCounter.Operator] // Default to 0 if not found
		txCountDiff := operatorTrxCounter.TrxCount - startTrxCount
		totalTrxCount += txCountDiff

		// Append result with the correct transaction difference
		operatorsTrxsCounts = append(operatorsTrxsCounts, types.OperatorTrxCounter{
			Operator: operatorTrxCounter.Operator,
			TrxCount: txCountDiff,
		})
	}

	return operatorsTrxsCounts, totalTrxCount
}

// IncrementOperatorTrxCount increments the transaction count for the given operator
// It is used to track the number of transactions per operator per block
// and is called from the AnteHandler to increment the transaction count
// It uses temporary storage to store the transaction count for the current block
func (k Keeper) IncrementOperatorTrxCount(ctx sdk.Context, operator string) error {
	// Retrieve the current block height
	blockHeight := ctx.BlockHeight()

	// Retrieve the current transaction count at the block height
	operatorsTrxsCount := k.GetOperatorTrxCount(ctx, blockHeight)
	if operatorsTrxsCount.BlockHeight == 0 {
		// Initialize the transaction count for the block height
		operatorsTrxsCount.BlockHeight = blockHeight
	}

	operators := k.GetParams(ctx).Operators

	// Increment the transaction count for the operator
	for _, _operator := range operators {
		if _operator == operator {
			// Increment the transaction count for the operator if it exists
			// Otherwise, add a new operator with a transaction count of 1
			for _, operatorTrxCounter := range operatorsTrxsCount.OperatorTrxCounters {
				if operatorTrxCounter.Operator == operator {
					operatorTrxCounter.TrxCount++
					operatorsTrxsCount.TrxCount++
					break
				}
			}
			operatorsTrxsCount.OperatorTrxCounters = append(operatorsTrxsCount.OperatorTrxCounters, &types.OperatorTrxCounter{
				Operator: operator,
				TrxCount: 1,
			})
			operatorsTrxsCount.TrxCount++
			break
		}
	}

	// Store the updated transaction count
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.MemTrxCounterKeyPrefix))
	b := k.cdc.MustMarshal(&operatorsTrxsCount)
	store.Set(types.MemTrxCounterKey(blockHeight), b)

	return nil
}

// GetOperatorTrxCount returns the transaction count for the given operator at the given block height
func (k Keeper) GetOperatorTrxCount(ctx sdk.Context, blockHeight int64) (operatorsTrxsCount types.OperatorsTrxsCount) {
	// Retrieve the transaction count at the given block height
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.MemTrxCounterKeyPrefix))
	b := store.Get(types.MemTrxCounterKey(blockHeight))
	if b == nil {
		return types.OperatorsTrxsCount{
			BlockHeight: blockHeight,
			TrxCount:    0,
		}
	}
	k.cdc.MustUnmarshal(b, &operatorsTrxsCount)
	return operatorsTrxsCount
}
