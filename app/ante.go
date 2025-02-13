package app

import (
	"chikku/x/egvmod/keeper"

	"github.com/cosmos/cosmos-sdk/x/auth/signing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// An ante handle that uses default BaseApp AnteHandler and uses it to count the number of transactions per block.
// NewAnteHandler creates a new AnteHandler
func NewAnteHandler(anteHandler sdk.AnteHandler, keeper *keeper.Keeper) (ah sdk.AnteHandler) {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		// Run the default BaseApp AnteHandler
		ah, err := anteHandler(ctx, tx, simulate)
		if err != nil {
			return ah, err
		}
		if !simulate {

			sigTx, _ := tx.(signing.SigVerifiableTx)

			signers, _ := sigTx.GetSigners()

			for _, signer := range signers {
				// convert to sdk.AccAddress
				accAddr := sdk.AccAddress(signer)

				params := keeper.GetParams(ctx)
				operatos := params.GetOperators()
				for _, operator := range operatos {
					if operator == (accAddr.String()) {
						keeper.IncrementOperatorTrxCount(ctx, operator)
						operatorTrxCount := keeper.GetOperatorTrxCount(ctx, ctx.BlockHeight())
						ctx.Logger().Info("Signers", "operatorTrxCount", operatorTrxCount)
						break
					}
				}
			}
		}

		//TODO: Increment the transaction counter
		//TODO: Assuming no multi sign trxs are enabled.
		return ah, nil
	}
}

//operatos=["chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv","chikku19tvrd6p8grkmxxus8r0df736kqn8mdfhefjx04","chikku15nq8krwa60ejpeedgtd3qq9ug29ruh6arelagl","chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4","chikku16gsjgtdslawyd0996qnm35hyv40deuc33huvnf","chikku1k7mtzld7089x4mu6w4wenl6ww5cfqd7dt3tljj","chikku1aae0ql3dzglukzks7zpsmwspaxz9xze5lunuvu","chikku1nrmmw6y25eckfl30pk3vutj27vysz49samxjmv","chikku1uwunwhkvt38l4u9tuve2kqejntgg0ye9updvuy","chikku1s5ssqjn799yse7gkjhfzr0nqzyat8kr0g9zley"]
// operatos=["chikku1ss6pjm244hy8nzk6zusypc5m9huneps40qvyvv","chikku19tvrd6p8grkmxxus8r0df736kqn8mdfhefjx04","chikku15nq8krwa60ejpeedgtd3qq9ug29ruh6arelagl","chikku1lxfdcxke9kr5hcyjex6g42q3lcpkcv7pttjes4","chikku16gsjgtdslawyd0996qnm35hyv40deuc33huvnf","chikku1k7mtzld7089x4mu6w4wenl6ww5cfqd7dt3tljj","chikku1aae0ql3dzglukzks7zpsmwspaxz9xze5lunuvu","chikku1nrmmw6y25eckfl30pk3vutj27vysz49samxjmv","chikku1uwunwhkvt38l4u9tuve2kqejntgg0ye9updvuy","chikku1s5ssqjn799yse7gkjhfzr0nqzyat8kr0g9zley"]
