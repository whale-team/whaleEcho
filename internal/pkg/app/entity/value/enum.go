package value

// MessageType represent message type
type MessageType int8

const (
	// Text represent text
	Text MessageType = iota + 1
	// File represent file
	File
)

// SysMsgType represent system message type
type SysMsgType int8

const (
	// CloseRoom represent room closed system message
	CloseRoom SysMsgType = iota + 1
)
