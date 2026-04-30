package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: en producción, validar origenes permitidos
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// ServeWS maneja el upgrade HTTP → WebSocket
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("firebase_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Error en upgrade WebSocket: %v", err)
			return
		}

		client := &Client{
			ID:   userID.(string),
			Hub:  hub,
			Conn: conn,
			Send: make(chan []byte, 256),
		}

		hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}

// ReadPump lee mensajes del WebSocket del cliente
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// Configurar conexión (solo funciona con *websocket.Conn real)
	if wsConn, ok := c.Conn.(*websocket.Conn); ok {
		wsConn.SetReadLimit(maxMessageSize)
		wsConn.SetReadDeadline(time.Now().Add(pongWait))
		wsConn.SetPongHandler(func(string) error {
			wsConn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	}

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error leyendo WS de %s: %v", c.ID, err)
			}
			break
		}

		// Decodificar el mensaje del cliente
		var payload MessagePayload
		if err := decodeMessage(message, &payload); err != nil {
			log.Printf("Mensaje inválido de %s: %v", c.ID, err)
			continue
		}

		// Asegurar que el sender sea el usuario autenticado
		payload.SenderID = c.ID

		// Enviar al hub para broadcast
		c.Hub.Broadcast <- &payload
	}
}

// WritePump envía mensajes al WebSocket del cliente
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if wsConn, okConn := c.Conn.(*websocket.Conn); okConn {
				wsConn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					wsConn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error escribiendo WS a %s: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			if wsConn, ok := c.Conn.(*websocket.Conn); ok {
				wsConn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
