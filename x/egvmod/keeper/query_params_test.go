package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "chikku/testutil/keeper"
	"chikku/x/egvmod/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.EgvmodKeeper(t)
	params := types.DefaultParams()
	require.NoError(t, keeper.SetParams(ctx, params))

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
