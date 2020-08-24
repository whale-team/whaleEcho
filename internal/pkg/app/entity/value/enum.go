package value

// MessageType represent message type
type MessageType int8

const (
	// Text represent text
	Text MessageType = iota + 1
	// File represent file
	File
)
