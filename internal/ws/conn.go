package ws

import "github.com/gorilla/websocket"

// WSConn abstrae la conexión WebSocket (facilita testing)
type WSConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// Verificar que *websocket.Conn implementa WSConn
var _ WSConn = (*websocket.Conn)(nil)
