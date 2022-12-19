package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	// get new id
	systemInfo, found := k.Keeper.GetSystemInfo(ctx) // k keeps the Keeper
	if !found {
		panic("SystemInfo not found") // it is ok to panic when there is no way to proceed due to something not a user error
	}
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	newGame := rules.New()
	storedGame := types.StoredGame{
		Index: newIndex,
		Board: newGame.String(),                 // new board state
		Turn:  rules.PieceStrings[newGame.Turn], // this returns "r" or "b" depending on rules
		Black: msg.Black,
		Red:   msg.Red,
	}

	// make sure the addresses black and red are valid
	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	k.Keeper.SetStoredGame(ctx, storedGame)

	// increase game id
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	return &types.MsgCreateGameResponse{}, nil
}
