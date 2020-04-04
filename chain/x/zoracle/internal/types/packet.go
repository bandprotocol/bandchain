package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
)

var _ channelexported.PacketDataI = OracleRequestPacketData{}
var _ channelexported.PacketDataI = OracleResponsePacketData{}

type OracleRequestPacketData struct {
	Data []byte `json:"data" yaml:"data"`
}

func NewOracleRequestPacketData(data []byte) OracleRequestPacketData {
	return OracleRequestPacketData{
		Data: data,
	}
}

func (o OracleRequestPacketData) String() string {
	return fmt.Sprintf(`OracleRequestPacketData:
	Data:                 %s`,
		o.Data,
	)
}

func (o OracleRequestPacketData) ValidateBasic() error {
	return nil
}

func (o OracleRequestPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(o))
}

func (o OracleRequestPacketData) GetTimeoutHeight() uint64 {
	return uint64(18446744073709551615)
}

func (o OracleRequestPacketData) Type() string {
	return "zoracle"
}

type OracleResponsePacketData struct {
	Data []byte `json:"data" yaml:"data"`
}

func NewOracleResponsePacketData(data []byte) OracleResponsePacketData {
	return OracleResponsePacketData{
		Data: data,
	}
}

func (o OracleResponsePacketData) String() string {
	return fmt.Sprintf(`OracleResponsePacketData:
	Data:                 %s`,
		o.Data,
	)
}

func (o OracleResponsePacketData) ValidateBasic() error {
	return nil
}

func (o OracleResponsePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(o))
}

func (o OracleResponsePacketData) GetTimeoutHeight() uint64 {
	return uint64(18446744073709551615)
}

func (o OracleResponsePacketData) Type() string {
	return "zoracle"
}
