package egvmod_test

import (
	"testing"

	keepertest "chikku/testutil/keeper"
	"chikku/testutil/nullify"
	egvmod "chikku/x/egvmod/module"
	"chikku/x/egvmod/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EgvmodKeeper(t)
	egvmod.InitGenesis(ctx, k, genesisState)
	got := egvmod.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
