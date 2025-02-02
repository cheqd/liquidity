package liquidity

import (
	"time"

	"github.com/cheqd/cosmos-sdk/telemetry"
	sdk "github.com/cheqd/cosmos-sdk/types"

	"github.com/gravity-devs/liquidity/x/liquidity/keeper"
	"github.com/gravity-devs/liquidity/x/liquidity/types"
)

// In the Begin blocker of the liquidity module,
// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.DeleteAndInitPoolBatches(ctx)
}

// In case of deposit, withdraw, and swap msgs, unlike other normal tx msgs,
// collect them in the liquidity pool batch and perform an execution once at the endblock to calculate and use the universal price.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.ExecutePoolBatches(ctx)
}
