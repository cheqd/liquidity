package testutil

import (
	"fmt"

	"github.com/cheqd/cosmos-sdk/baseapp"
	"github.com/cheqd/cosmos-sdk/client"
	"github.com/cheqd/cosmos-sdk/client/flags"
	servertypes "github.com/cheqd/cosmos-sdk/server/types"
	"github.com/cheqd/cosmos-sdk/simapp"
	"github.com/cheqd/cosmos-sdk/simapp/params"
	storetypes "github.com/cheqd/cosmos-sdk/store/types"
	"github.com/cheqd/cosmos-sdk/testutil"
	clitestutil "github.com/cheqd/cosmos-sdk/testutil/cli"
	"github.com/cheqd/cosmos-sdk/testutil/network"
	sdk "github.com/cheqd/cosmos-sdk/types"
	govcli "github.com/cheqd/cosmos-sdk/x/gov/client/cli"
	paramscli "github.com/cheqd/cosmos-sdk/x/params/client/cli"

	liquidityapp "github.com/gravity-devs/liquidity/app"
	liquiditycli "github.com/gravity-devs/liquidity/x/liquidity/client/cli"

	dbm "github.com/tendermint/tm-db"
)

// NewConfig returns config that defines the necessary testing requirements
// used to bootstrap and start an in-process local testing network.
func NewConfig(dbm *dbm.MemDB) network.Config {
	encCfg := simapp.MakeTestEncodingConfig()

	cfg := network.DefaultConfig()
	cfg.AppConstructor = NewAppConstructor(encCfg, dbm)                    // the ABCI application constructor
	cfg.GenesisState = liquidityapp.ModuleBasics.DefaultGenesis(cfg.Codec) // liquidity genesis state to provide
	return cfg
}

// NewAppConstructor returns a new network AppConstructor.
func NewAppConstructor(encodingCfg params.EncodingConfig, db *dbm.MemDB) network.AppConstructor {
	return func(val network.Validator) servertypes.Application {
		return liquidityapp.NewLiquidityApp(
			val.Ctx.Logger, db, nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			liquidityapp.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

// MsgCreatePoolExec creates a transaction for creating liquidity pool.
func MsgCreatePoolExec(clientCtx client.Context, from, poolID, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewCreatePoolCmd(), args)
}

// MsgDepositWithinBatchExec creates a transaction to deposit new amounts to the pool.
func MsgDepositWithinBatchExec(clientCtx client.Context, from, poolID, depositCoins string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewDepositWithinBatchCmd(), args)
}

// MsgWithdrawWithinBatchExec creates a transaction to withraw pool coin amount from the pool.
func MsgWithdrawWithinBatchExec(clientCtx client.Context, from, poolID, poolCoin string,
	extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		poolCoin,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewWithdrawWithinBatchCmd(), args)
}

// MsgSwapWithinBatchExec creates a transaction to swap coins in the pool.
func MsgSwapWithinBatchExec(clientCtx client.Context, from, poolID, swapTypeID,
	offerCoin, demandCoinDenom, orderPrice, swapFeeRate string, extraArgs ...string) (testutil.BufferWriter, error) {

	args := append([]string{
		poolID,
		swapTypeID,
		offerCoin,
		demandCoinDenom,
		orderPrice,
		swapFeeRate,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewSwapWithinBatchCmd(), args)
}

// MsgParamChangeProposalExec creates a transaction for submitting param change proposal
func MsgParamChangeProposalExec(clientCtx client.Context, from string, file string) (testutil.BufferWriter, error) {

	args := append([]string{
		file,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	paramChangeCmd := paramscli.NewSubmitParamChangeProposalTxCmd()
	flags.AddTxFlagsToCmd(paramChangeCmd)

	return clitestutil.ExecTestCLICmd(clientCtx, paramChangeCmd, args)
}

// MsgVote votes for a proposal
func MsgVote(clientCtx client.Context, from, id, vote string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := append([]string{
		id,
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdWeightedVote(), args)
}
