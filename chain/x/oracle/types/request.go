package types

var (
	_ RequestSpec = &MsgRequestData{}
	_ RequestSpec = &OracleRequestPacketData{}
)

type RequestSpec interface {
	GetOracleScriptID() OracleScriptID
	GetCalldata() []byte
	GetAskCount() int64
	GetMinCount() int64
	GetClientID() string
}
