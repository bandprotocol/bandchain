package types

var (
	_ RequestSpec = &MsgRequestData{}
	_ RequestSpec = &OracleRequestPacketData{}
)

// RequestSpec captures the essence of what it means to be a request-making object.
type RequestSpec interface {
	GetOracleScriptID() OracleScriptID
	GetCalldata() []byte
	GetAskCount() uint64
	GetMinCount() uint64
	GetClientID() string
}
