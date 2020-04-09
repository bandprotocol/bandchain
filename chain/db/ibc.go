package db

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

func (b *BandDB) getChainID(channelID, channelPort string) (string, error) {
	channel, found := b.IBCKeeper.ChannelKeeper.GetChannel(
		b.ctx,
		channelPort,
		channelID,
	)
	if !found {
		return "", sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot find channel")
	}
	connection, found := b.IBCKeeper.ConnectionKeeper.GetConnection(b.ctx, channel.ConnectionHops[0])
	if !found {
		return "", sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot find connection")
	}

	client, found := b.IBCKeeper.ClientKeeper.GetClientState(b.ctx, connection.GetClientID())
	if !found {
		return "", sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot find client")
	}
	return client.GetChainID(), nil
}

func (b *BandDB) handleMsgPacket(
	txHash []byte,
	msg channel.MsgPacket,
	events map[string]string,
) error {
	packetType := ""
	jsonMap := make(map[string]interface{})

	err := json.Unmarshal(msg.GetData(), &jsonMap)
	if err != nil {
		return err
	}
	extra := make(map[string]interface{})

	var requestData oracle.OracleRequestPacketData
	if err := oracle.ModuleCdc.UnmarshalJSON(msg.GetData(), &requestData); err == nil {
		packetType = "ORACLE REQUEST"
		id, err := strconv.ParseInt(events[oracle.EventTypeRequest+"."+oracle.AttributeKeyID], 10, 64)
		if err != nil {
			return err
		}
		calldata, err := hex.DecodeString(requestData.Calldata)
		if err != nil {
			return err
		}
		err = b.AddNewRequest(
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
		if err != nil {
			return err
		}

		extra["requestID"] = id
	}
	if packetType == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized oracle package type: %T", msg.Packet)
	}

	chainID, err := b.getChainID(msg.GetDestChannel(), msg.GetDestPort())
	if err != nil {
		return err
	}

	jsonMap["extra"] = extra
	rawJson, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	return b.tx.Create(&Packet{
		Type:        packetType,
		Sequence:    msg.GetSequence(),
		MyChannel:   msg.GetDestChannel(),
		MyPort:      msg.GetDestPort(),
		YourChainID: chainID,
		YourChannel: msg.GetSourceChannel(),
		YourPort:    msg.GetSourcePort(),
		BlockHeight: b.ctx.BlockHeight(),
		IsIncoming:  true,
		Detail:      rawJson,
	}).Error
}
