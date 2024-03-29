package keeper

import (
	"fmt"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	// differentiate between players. Players pay in their first move
	if storedGame.MoveCount == 0 {
		// Black plays first
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		// if address aquired, then escrow the money
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrBlackCannotPay.Error())
		}
	} else if storedGame.MoveCount == 1 {
		// Red plays second
		red, err := storedGame.GetRedAddress()
		if err != nil {
			panic(err.Error())
		}
		// if address aquired, then escrow the money
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrRedCannotPay.Error())
		}
	}
	return nil
}
func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {
	// get winner address
	winnerAddress, found, err := storedGame.GetWinnerAddress()
	if err != nil {
		panic(err.Error())
	}
	if !found {
		panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
	}
	// determine amount to pay
	winnings := storedGame.GetWagerCoin()
	if storedGame.MoveCount == 0 {
		panic(types.ErrNothingToPay.Error())
	} else if 1 < storedGame.MoveCount {
		winnings = winnings.Add(winnings)
	}
	// pay the winnings
	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winnerAddress, sdk.NewCoins(winnings))
	if err != nil {
		panic(fmt.Sprintf(types.ErrCannotPayWinnings.Error(), err.Error()))
	}
}
func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {
	if storedGame.MoveCount == 1 {
		// Refund the black player
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), err.Error()))
		}
	} else if storedGame.MoveCount == 0 {
		// Do nothing
	} else {
		// TODO Implement a draw mechanism.
		panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
	}
}
