package socketserver

func handlePingMessage(wsConnection *WsConnection) {
	ServerManager().WriteConnection(wsConnection, SocketResponsePing{Ping: "pong"})
}

// SocketResponsePing is the ping response
type SocketResponsePing struct {
	Ping string `json:"ping"`
}
