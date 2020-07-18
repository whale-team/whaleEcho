package value

type MessageType int8

const (
	Text MessageType = iota + 1
	File
	Join
	Leave
)
