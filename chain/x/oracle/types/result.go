package types

// Result is a convenience struct that keeps both request and response packets of a request.
type Result struct {
	RequestPacketData  OracleRequestPacketData
	ResponsePacketData OracleResponsePacketData
}
