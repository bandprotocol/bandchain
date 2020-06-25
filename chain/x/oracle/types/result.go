package types

// Result is a convenience struct that keeps both request and respose packets of a request.
type Result struct {
	RequestPacketData  OracleRequestPacketData
	ResponsePacketData OracleResponsePacketData
}
