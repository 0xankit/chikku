package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "chikku/testutil/keeper"
	"chikku/x/egvmod/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.EgvmodKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
