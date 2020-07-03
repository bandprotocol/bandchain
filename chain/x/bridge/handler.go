package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgUpdateChainID:
			return handleMsgUpdateChainID(ctx, k, msg)
		case MsgUpdateValidators:
			return handleMsgUpdateValidators(ctx, k, msg)
		case MsgRelay:
			return handleMsgRelay(ctx, k, msg)
		case MsgVerifyProof:
			return handleMsgVerifyProof(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgUpdateChainID(ctx sdk.Context, k Keeper, m MsgUpdateChainID) (*sdk.Result, error) {
	// TODO: Add validate only owner
	k.SetChainID(ctx, m.ChainID)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdateValidators(ctx sdk.Context, k Keeper, m MsgUpdateValidators) (*sdk.Result, error) {
	// TODO: Add validate only owner
	k.UpdateValidators(ctx, m.Validators)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRelay(ctx sdk.Context, k Keeper, m MsgRelay) (*sdk.Result, error) {
	err := k.Relay(ctx, m.SignedHeader)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgVerifyProof(ctx sdk.Context, k Keeper, m MsgVerifyProof) (*sdk.Result, error) {
	err := k.VerifyProof(ctx, m.Height, m.Proof, m.RequestPacket, m.ResponsePacket)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
