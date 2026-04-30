package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// MessagePayload es lo que viaja por el WebSocket
type MessagePayload struct {
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	Timestamp      string `json:"timestamp"`
}

// Client representa una conexión WebSocket activa
type Client struct {
	ID   string
	Hub  *Hub
	Conn WSConn // interfaz para facilitar testing
	Send chan []byte
}

// Hub mantiene el registro de todos los clientes conectados
type Hub struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *MessagePayload
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *MessagePayload),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Cliente conectado: %s (total: %d)", client.ID, len(h.clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				close(client.Send)
				delete(h.clients, client.ID)
			}
			h.mu.Unlock()
			log.Printf("Cliente desconectado: %s", client.ID)

		case msg := <-h.Broadcast:
			data := encodeMessage(msg)
			h.mu.RLock()
			for id, client := range h.clients {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.clients, id)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// encodeMessage serializa un MessagePayload a JSON bytes
func encodeMessage(msg *MessagePayload) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error serializando mensaje: %v", err)
		return []byte("{}")
	}
	return data
}
