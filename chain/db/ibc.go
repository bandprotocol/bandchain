package db

import (
	"encoding/hex"
	"strconv"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

func (b *BandDB) handleMsgPacket(
	txHash []byte,
	msg channel.MsgPacket,
	events map[string]string,
) error {
	var requestData oracle.OracleRequestPacketData
	if err := oracle.ModuleCdc.UnmarshalJSON(msg.GetData(), &requestData); err == nil {
		id, err := strconv.ParseInt(events[oracle.EventTypeRequest+"."+oracle.AttributeKeyID], 10, 64)
		if err != nil {
			return err
		}
		calldata, err := hex.DecodeString(requestData.Calldata)
		if err != nil {
			return err
		}
		return b.AddNewRequest(
			id,
			int64(requestData.OracleScriptID),
			calldata,
			requestData.SufficientValidatorCount,
			requestData.Expiration,
			"Pending",
			msg.Signer.String(),
			requestData.ClientID,
			txHash,
			nil,
		)
	}
	return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized oracle package type: %T", msg.Packet)
}
