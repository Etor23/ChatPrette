package ws

import "encoding/json"

// decodeMessage deserializa JSON bytes a un MessagePayload
func decodeMessage(data []byte, payload *MessagePayload) error {
	return json.Unmarshal(data, payload)
}
