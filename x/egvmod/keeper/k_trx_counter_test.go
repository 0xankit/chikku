package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "chikku/testutil/keeper"
)

func TestIncrementOperatorTrxCount(t *testing.T) {
	k, ctx := keepertest.EgvmodKeeper(t) // Assume setupKeeper initializes a keeper with mock storage

	operator := "chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv"
	// add block height to context
	ctx = ctx.WithBlockHeight(1)
	err := k.IncrementOperatorTrxCount(ctx, operator)
	require.NoError(t, err)

	count := k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	require.Equal(t, int64(1), count.TrxCount)

	// Increment again
	err = k.IncrementOperatorTrxCount(ctx, operator)
	require.NoError(t, err)

	count = k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	require.Equal(t, int64(2), count.TrxCount)
}

func TestSetTrxCount(t *testing.T) {
	k, ctx := keepertest.EgvmodKeeper(t)
	// add block height to context
	ctx = ctx.WithBlockHeight(1)

	operator := "chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv"
	_ = k.IncrementOperatorTrxCount(ctx, operator)

	trxCount := k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	require.Equal(t, int64(1), trxCount.TrxCount)

	// Persist transactions
	k.SetTrxCount(ctx, trxCount)

	// Check if stored correctly
	storedCount := k.GetTrxCount(ctx, ctx.BlockHeight())
	require.Equal(t, int64(1), storedCount.TrxCount)

	// Check if temporary data is cleared
	tempCount := k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	require.Equal(t, int64(0), tempCount.TrxCount)
}

func TestGetTrxCountsInRange(t *testing.T) {
	k, ctx := keepertest.EgvmodKeeper(t)

	// Start at block height 1
	ctx = ctx.WithBlockHeight(1)

	// Increment transactions for multiple operators
	_ = k.IncrementOperatorTrxCount(ctx, "chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv") // +1
	_ = k.IncrementOperatorTrxCount(ctx, "chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv") // +1
	_ = k.IncrementOperatorTrxCount(ctx, "chikku19tvrd6p8grkmxxus8r0df736kqn8mdfhefjx04") // +1
	_ = k.IncrementOperatorTrxCount(ctx, "chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4") // +1

	// Persist transactions for block 1
	trxCount := k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	k.SetTrxCount(ctx, trxCount)

	startHeight := ctx.BlockHeight()           // Save this block height
	ctx = ctx.WithBlockHeight(startHeight + 1) // Move to next block

	// Increment transaction count in block 2
	_ = k.IncrementOperatorTrxCount(ctx, "chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4") // +1

	// Persist transactions for block 2
	newTrxCount := k.GetOperatorTrxCount(ctx, ctx.BlockHeight())
	k.SetTrxCount(ctx, newTrxCount)

	// Get transaction counts between `startHeight` (block 1) and current height (block 2)
	operatorsTrxsCounts, totalTrx := k.GetTrxCountsInRange(ctx, startHeight, ctx.BlockHeight())

	// Validate total transaction count
	require.Equal(t, int64(5), totalTrx, "Total transaction count mismatch")

	// Validate per-operator transaction count
	expectedCounts := map[string]int64{
		"chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv": 2,
		"chikku19tvrd6p8grkmxxus8r0df736kqn8mdfhefjx04": 1,
		"chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4": 2, // 1 from first block + 1 from second block
	}

	require.Len(t, operatorsTrxsCounts, 3, "Unexpected number of operators in transaction range")

	for _, operatorTrx := range operatorsTrxsCounts {
		expectedCount, exists := expectedCounts[operatorTrx.Operator]
		require.True(t, exists, "Unexpected operator in transaction range: %s", operatorTrx.Operator)
		require.Equal(t, expectedCount, operatorTrx.TrxCount, "Incorrect transaction count for operator: %s", operatorTrx.Operator)
	}
}
