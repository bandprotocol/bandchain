package types

// Result is a convenience struct that keeps both request and response packets of a request.
type Result struct {
	RequestPacketData  OracleRequestPacketData  `json:"request_packet_data"`
	ResponsePacketData OracleResponsePacketData `json:"response_packet_data"`
}

// NewResult creates a new Result instance.
func NewResult(req OracleRequestPacketData, res OracleResponsePacketData) Result {
	return Result{
		RequestPacketData:  req,
		ResponsePacketData: res,
	}
}
