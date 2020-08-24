package entity

// Message represent websocket message received from client
type Message struct {
	Data       []byte
	RoomUID    string
	SenderName string
}
