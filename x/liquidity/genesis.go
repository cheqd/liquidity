package liquidity

import (
	sdk "github.com/cheqd/cosmos-sdk/types"

	"github.com/gravity-devs/liquidity/x/liquidity/keeper"
	"github.com/gravity-devs/liquidity/x/liquidity/types"
)

// InitGenesis new liquidity genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.InitGenesis(ctx, data)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return keeper.ExportGenesis(ctx)
}
