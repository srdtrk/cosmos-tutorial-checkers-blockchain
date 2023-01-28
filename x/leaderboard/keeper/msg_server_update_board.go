package keeper

import (
	"context"

	"github.com/alice/checkers/x/leaderboard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateBoard(goCtx context.Context, msg *types.MsgUpdateBoard) (*types.MsgUpdateBoardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgUpdateBoardResponse{}, nil
}
