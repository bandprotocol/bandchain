package db

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
	events map[string][]string,
) error {
	packetType := ""
	jsonMap := make(map[string]interface{})
	var requestPacket oracle.OracleRequestPacketData
	if err := oracle.ModuleCdc.UnmarshalJSON(msg.GetData(), &requestPacket); err == nil {
		packetType = "ORACLE REQUEST"
		rawBytes, err := json.Marshal(requestPacket)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rawBytes, &jsonMap)
		if err != nil {
			return err
		}
		ids := events[otypes.EventTypeRequest+"."+otypes.AttributeKeyID]
		if len(ids) != 1 {
			return errors.New("Cannot find request id")
		}
		id, err := strconv.ParseInt(ids[0], 10, 64)
		if err != nil {
			return err
		}
		request, err := b.OracleKeeper.GetRequest(b.ctx, otypes.RequestID(id))
		if err != nil {
			return err
		}
		err = b.AddNewRequest(
			id,
			int64(requestPacket.OracleScriptID),
			requestPacket.Calldata,
			requestPacket.MinCount,
			request.RequestHeight+20, // TODO: REMOVE THIS. HACK!!!!
			"Pending",
			msg.Signer.String(),
			requestPacket.ClientID,
			txHash,
			nil,
			events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyExternalID],
			events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyDataSourceID],
			events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyCalldata],
		)
		if err != nil {
			return err
		}

		oracleScript, err := b.OracleKeeper.GetOracleScript(b.ctx, requestPacket.OracleScriptID)
		if err != nil {
			return err
		}

		jsonMap["type"] = "oracle/OracleRequestPacketData"
		jsonMap["request_id"] = id
		jsonMap["oracle_script_name"] = oracleScript.Name
		jsonMap["schema"] = oracleScript.Schema
	}
	if packetType == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized oracle package type: %T", msg.Packet)
	}

	chainID, err := b.getChainID(msg.GetDestChannel(), msg.GetDestPort())
	if err != nil {
		return err
	}

	rawJson, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	isIncoming := true
	return b.tx.Create(&Packet{
		Type:        packetType,
		Sequence:    msg.GetSequence(),
		MyChannel:   msg.GetDestChannel(),
		MyPort:      msg.GetDestPort(),
		YourChainID: chainID,
		YourChannel: msg.GetSourceChannel(),
		YourPort:    msg.GetSourcePort(),
		BlockHeight: b.ctx.BlockHeight(),
		IsIncoming:  &isIncoming,
		Detail:      rawJson,
	}).Error
}
