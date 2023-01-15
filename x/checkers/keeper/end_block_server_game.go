package keeper

import (
	"context"
	"fmt"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ForfeitExpiredGames(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	opponents := map[string]string{
		rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER],
		rules.PieceStrings[rules.RED_PLAYER]:   rules.PieceStrings[rules.BLACK_PLAYER],
	}

	systemInfo, found := k.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	gameIndex := systemInfo.FifoHeadIndex
	var storedGame types.StoredGame

	for {
		// break if FIFO reached an end
		if gameIndex == types.NoFifoIndex {
			break
		}

		// Fetch the game and its deadline
		storedGame, found = k.GetStoredGame(ctx, gameIndex)
		if !found {
			panic("Fifo head game not found " + systemInfo.FifoHeadIndex)
		}
		deadline, err := storedGame.GetDeadlineAsTime()
		if err != nil {
			panic(err)
		}

		// Test for expiration
		if deadline.Before(ctx.BlockTime()) {
			// remove it from FIFO
			k.RemoveFromFifo(ctx, &storedGame, &systemInfo)

			// Determine if the game is worth keeping (i.e. whether or not we should pretend the game never existed)
			// if so, then determine the winner, which is the opponent of the player that didn't make their move before the deadline
			lastBoard := storedGame.Board
			if storedGame.MoveCount <= 1 {
				// No point in keeping a game that was never really played
				k.RemoveStoredGame(ctx, gameIndex)
				// the game was never really played. Refund the wager of the player who started the game.
				if storedGame.MoveCount == 1 {
					k.MustRefundWager(ctx, &storedGame)
				}
			} else {
				storedGame.Winner, found = opponents[storedGame.Turn]
				if !found {
					panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
				}
				storedGame.Board = ""
				// Pay the winnings of the player who won the game
				k.MustPayWinnings(ctx, &storedGame)
				k.SetStoredGame(ctx, storedGame)
			}
			// emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.GameForfeitedEventType,
					sdk.NewAttribute(types.GameForfeitedEventGameIndex, gameIndex),
					sdk.NewAttribute(types.GameForfeitedEventWinner, storedGame.Winner),
					sdk.NewAttribute(types.GameForfeitedEventBoard, lastBoard),
				),
			)
			// Move in FIFO
			gameIndex = systemInfo.FifoHeadIndex
		} else {
			// All other games after are active anyway
			break
		}
	}

	k.SetSystemInfo(ctx, systemInfo)
}
