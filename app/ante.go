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
		ctx.Logger().Error("NewAnteHandler called")
		// Run the default BaseApp AnteHandler
		ah, err := anteHandler(ctx, tx, simulate)
		if err != nil {
			return ah, err
		}
		if !simulate {

			sigTx, _ := tx.(signing.SigVerifiableTx)

			signers, _ := sigTx.GetSigners()

			for _, signer := range signers {
				ctx.Logger().Error("Signers", "sig", signer)
				// convert to sdk.AccAddress
				accAddr := sdk.AccAddress(signer)
				ctx.Logger().Error("Signers", "accAddr", accAddr)

				params := keeper.GetParams(ctx)
				operatos := params.GetOperators()
				ctx.Logger().Error("Signers", "operatos", operatos)
				for _, operator := range operatos {
					if operator == (accAddr.String()) {
						ctx.Logger().Error("####################### Signers", "operator", operator)
						keeper.IncrementOperatorTrxCount(ctx, operator)
						// os.Exit(1)
						operatorTrxCount := keeper.GetOperatorTrxCount(ctx, ctx.BlockHeight())
						ctx.Logger().Error("Signers", "operatorTrxCount", operatorTrxCount)
						break
					}
					ctx.Logger().Error("Signers", "operator", operator, "accAddr", accAddr)
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
