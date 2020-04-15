package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OracleRequestPacketData encodes an oracle request sent from other blockchains to BandChain.
type OracleRequestPacketData struct {
	// ClientID is the unique identifier of this oracle request, as specified by the client.
	// This same unique ID will be sent back to the requester with the oracle response.
	ClientID string `json:"client_id" yaml:"client_id"`
	// OracleScriptID is the unique identifier of the oracle script to be executed.
	OracleScriptID OracleScriptID `json:"oracle_script_id" yaml:"oracle_script_id"`
	// Calldata is the hex-encoded of the calldata bytes available for oracle execution
	// during both preparation and execution phases.
	Calldata string `json:"calldata" yaml:"calldata"`
	// AskCount is the number of validators that are requested to respond to this oracle
	// request. Higher value means more security, at a higher gas cost.
	AskCount int64 `json:"ask_count" yaml:"ask_count"`
	// MinCount is the minimum number of validators necessary for the request to proceed to
	// the execution phase. Higher value means more security, at the cost of liveness.
	MinCount int64 `json:"min_count" yaml:"min_count"`
}

// OracleResponsePacketData encodes an oracle response from BandChain to the requester.
type OracleResponsePacketData struct {
	// ClientID is the unique identifier matched with that of the oracle request packet.
	ClientID string `json:"client_id" yaml:"client_id"`
	// RequestID is BandChain's unique identifier for this oracle request.
	// TODO: This is not actually needed, but is here to simplify DB. Should remove.
	RequestID RequestID `json:"request_id" yaml:"request_id"`
	// AnsCount is the number of validators among to the asked validators that actually
	// responded to this oracle request prior to this oracle request being resolved.
	AnsCount int64 `json:"ans_count" yaml:"ans_count"`
	// PrepareTime is the UNIX epoch time at which the request was sent to BandChain.
	PrepareTime int64 `json:"prepare_time" yaml:"prepare_time"`
	// ResolveTime is the UNIX epoch time at which the request was resolved to the final result.
	ResolveTime int64 `json:"resolve_time" yaml:"resolve_time"`
	// ResolveStatus is the status of this oracle request, which can be OK, ERROR, or EXPIRED.
	ResolveStatus ResolveStatus `json:"resolve_status" yaml:"resolve_status"`
	// Result is the hex-encoded of the final aggregated value only available if status if OK.
	Result string `json:"result" yaml:"result"`
}

// NewOracleRequestPacketData creates a new OracleRequestPacketData instance.
func NewOracleRequestPacketData(
	clientID string, oracleScriptID OracleScriptID, calldata string,
	askCount int64, minCount int64,
) OracleRequestPacketData {
	return OracleRequestPacketData{
		ClientID:       clientID,
		OracleScriptID: oracleScriptID,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
	}
}

func (p OracleRequestPacketData) String() string {
	return fmt.Sprintf(`OracleRequestPacketData:
    ClientID:       %s
    OracleScriptID: %d
    Calldata:       %s
    AskCount:       %d
    MinCount:       %d`,
		p.ClientID,
		p.OracleScriptID,
		p.Calldata,
		p.AskCount,
		p.MinCount,
	)
}

func (p OracleRequestPacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleRequestPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

func NewOracleResponsePacketData(
	clientID string, requestID RequestID, ansCount int64, prepareTime int64, resolveTime int64,
	resolveStatus ResolveStatus, result string,
) OracleResponsePacketData {
	return OracleResponsePacketData{
		ClientID:      clientID,
		RequestID:     requestID,
		AnsCount:      ansCount,
		PrepareTime:   prepareTime,
		ResolveTime:   resolveTime,
		ResolveStatus: resolveStatus,
		Result:        result,
	}
}

func (p OracleResponsePacketData) String() string {
	return fmt.Sprintf(`OracleResponsePacketData:
	ClientID: %s
	RequestID: %d
	AnsCount: %d
	PrepareTime: %d
	ResolveTime: %d
	ResolveStatus: %d
	Result: %s`,
		p.ClientID,
		p.RequestID,
		p.AnsCount,
		p.PrepareTime,
		p.ResolveTime,
		p.ResolveStatus,
		p.Result,
	)
}

func (p OracleResponsePacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleResponsePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}
