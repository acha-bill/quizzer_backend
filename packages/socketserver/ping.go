package socketserver

func handlePingMessage(wsConnection *WsConnection) {
	ServerManager().WriteConnection(wsConnection, SocketResponsePing{Type: pingResponseType, Ping: "pong"})
}

const pingResponseType = "pong"

// SocketResponsePing is the ping response
type SocketResponsePing struct {
	Type string `json:"type"`
	Ping string `json:"ping"`
}
