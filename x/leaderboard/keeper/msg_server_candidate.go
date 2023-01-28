package keeper

import (
	"context"
	"errors"

	"github.com/alice/checkers/x/leaderboard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
)

func (k msgServer) SendCandidate(goCtx context.Context, msg *types.MsgSendCandidate) (*types.MsgSendCandidateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Construct the packet
	var packet types.CandidatePacketData

	allPlayerInfo := k.GetAllPlayerInfo(ctx)

	found_in_player_list := false
	for i := range allPlayerInfo {
		if allPlayerInfo[i].Index == msg.Creator {
			packet.PlayerInfo = &allPlayerInfo[i]
			found_in_player_list = true
			break
		}
	}

	if !found_in_player_list {
		return nil, errors.New("player not found")
	}

	// Transmit the packet
	err := k.TransmitCandidatePacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendCandidateResponse{}, nil
}
