package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"chikku/x/egvmod/keeper"
	"chikku/x/egvmod/types"
)

func EgvmodKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		nil,
		nil,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	operators := []string{"chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv",
		"chikku19tvrd6p8grkmxxus8r0df736kqn8mdfhefjx04",
		"chikku15nq8krwa60ejpeedgtd3qq9ug29ruh6arelagl",
		"chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4",
		"chikku16gsjgtdslawyd0996qnm35hyv40deuc33huvnf",
		"chikku1k7mtzld7089x4mu6w4wenl6ww5cfqd7dt3tljj",
		"chikku1aae0ql3dzglukzks7zpsmwspaxz9xze5lunuvu",
		"chikku1nrmmw6y25eckfl30pk3vutj27vysz49samxjmv",
		"chikku1uwunwhkvt38l4u9tuve2kqejntgg0ye9updvuy",
		"chikku1s5ssqjn799yse7gkjhfzr0nqzyat8kr0g9zley",
	}

	params := types.DefaultParams()
	params.Operators = operators
	// Initialize params
	if err := k.SetParams(ctx, params); err != nil {
		panic(err)
	}

	return k, ctx
}
