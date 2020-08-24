package socketserver

func handlePingMessage(wsConnection *WsConnection) {
	Manager().WriteConnection(wsConnection, SocketResponsePing{Ping: "pong"})
}

// SocketResponsePing is the ping response
type SocketResponsePing struct {
	Ping string `json:"ping"`
}
