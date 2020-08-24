package socketserver

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/acha-bill/quizzer_backend/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	upgrader               = websocket.Upgrader{}
	ErrWritingToConnection = errors.New("error writing to connection")
	ErrUnknownMessageType  = errors.New("unknown message type")
	msgTypeMap             map[string]interface{}
	once                   sync.Once
	mutex                  sync.Mutex
	manager                *WsManager
)

const (
	MessageTypeAuth   = "auth"
	MessageTypePing   = "ping"
	MessageTypeAnswer = "answer"
)

func init() {
	msgTypeMap = make(map[string]interface{})
	msgTypeMap[MessageTypeAuth] = SocketMessageAuth{}
	msgTypeMap[MessageTypePing] = nil
	msgTypeMap[MessageTypeAnswer] = SocketMessageAnswer{}
}

// WsContext is the context of a socket connection
type WsContext struct {
	Ready bool
	User  *models.User
}

// WsConnection represents the websocket connection.
type WsConnection struct {
	Socket  *websocket.Conn
	Context *WsContext
}

// WsManager is the manager of connections
type WsManager struct {
	connections map[*websocket.Conn]*WsConnection
}

// ServerManager returns the wsmanager instance
func ServerManager() *WsManager {
	once.Do(func() {
		manager = &WsManager{
			connections: make(map[*websocket.Conn]*WsConnection),
		}
	})
	return manager
}

// AddConnection adds a new socket connection
func (mgr *WsManager) AddConnection(wsCon *WsConnection) {
	mutex.Lock()
	mgr.connections[wsCon.Socket] = wsCon
	mutex.Unlock()
}

// RemoveConnection closes and removes the connection from the manager
func (mgr *WsManager) RemoveConnection(conn *websocket.Conn) {
	mutex.Lock()
	delete(mgr.connections, conn)
	mutex.Unlock()
}

// CloseConnection closes the connection
func (mgr *WsManager) CloseConnection(conn *websocket.Conn) {
	conn.Close()
}

// Length returns the number of active connections
func (mgr *WsManager) Length() int {
	mutex.Lock()
	l := len(mgr.connections)
	mutex.Unlock()
	return l
}

// Get gets the specified connection
func (mgr *WsManager) Get(conn *websocket.Conn) *WsConnection {
	return mgr.connections[conn]
}

// WriteConnection writes the JSON serialized form of the data to the connection
func (mgr *WsManager) WriteConnection(conn *WsConnection, data interface{}) {
	if err := conn.Socket.WriteJSON(data); err != nil {
		log.Errorf("%w", ErrWritingToConnection)
	}
}

// Listen listens for connections and reads messages
func Listen(ctx echo.Context) error {
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		// no need to close as its already closing
		ServerManager().RemoveConnection(conn)
		conn.Close()
		return nil
	})

	wsConn := &WsConnection{
		Socket:  conn,
		Context: &WsContext{Ready: false},
	}
	ServerManager().AddConnection(wsConn)

	for {
		// Read
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("%v", err)
			_ = conn.WriteJSON(SocketResponseError{Error: err.Error()})

		} else {
			go handleRead(bytes, conn)
		}
	}
}

func handleRead(bytes []byte, conn *websocket.Conn) {
	var msg SocketMessage
	if err := json.Unmarshal(bytes, &msg); err != nil {
		log.Errorf("%v", err)
		_ = conn.WriteJSON(SocketResponseError{Error: err.Error()})
		return
	}

	target, ok := msgTypeMap[msg.Type]
	if !ok {
		_ = conn.WriteJSON(SocketResponseError{Error: ErrUnknownMessageType.Error()})
		return
	}

	// first read must be auth
	wsConnection := ServerManager().Get(conn)
	if !wsConnection.Context.Ready && msg.Type != MessageTypeAuth && msg.Type != MessageTypePing {
		_ = conn.WriteJSON(SocketResponseError{Error: ErrSocketNotAuthenticated.Error()})
		return
	}

	// handlers
	switch msg.Type {
	case MessageTypePing:
		handlePingMessage(wsConnection)
	case MessageTypeAuth:
		authMsg := target.(SocketMessageAuth)
		handleAuthMessage(wsConnection, authMsg)
	case MessageTypeAnswer:
		answerMsg := target.(SocketMessageAnswer)
		handleAnswerMessage(wsConnection, answerMsg)
	}
}

type SocketResponseError struct {
	Error string `json:"error,omitempty"`
}

type SocketMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
