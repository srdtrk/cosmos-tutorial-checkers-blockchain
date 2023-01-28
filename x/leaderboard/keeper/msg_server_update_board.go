package keeper

import (
	"context"

	"github.com/alice/checkers/x/leaderboard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateBoard(goCtx context.Context, msg *types.MsgUpdateBoard) (*types.MsgUpdateBoardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	playerInfoList := k.GetAllPlayerInfo(ctx)
	k.updateBoard(ctx, playerInfoList)

	return &types.MsgUpdateBoardResponse{}, nil
}
