package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AddressPrefix = "chikku"
)

func init() {
	// Make sure that addresses the same prefix in blockchain and tests/simulation
	sdk.GetConfig().SetBech32PrefixForAccount(AddressPrefix, "pub")
	sdk.SetBaseDenom("egv")
}

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func Address() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(AccAddress())
}

func Coin(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

func Operators(addresses []string) []sdk.AccAddress {
	operators := make([]sdk.AccAddress, len(addresses))
	for i, addr := range addresses {
		operators[i] = sdk.MustAccAddressFromBech32(addr)
	}
	return operators
}
