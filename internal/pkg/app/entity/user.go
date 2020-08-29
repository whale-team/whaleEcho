package entity

// BinaryWriterCloser define write close method
type BinaryWriterCloser interface {
	WriteBinary([]byte) error
	Close() error
}

// User represent mapping between client connection and its identitiy
type User struct {
	conn    BinaryWriterCloser
	Name    string
	UID     string
	RoomUID string
}

// BindConn binding websocket connection
func (u *User) BindConn(connCtx BinaryWriterCloser) {
	u.conn = connCtx
}

// CloseConn close the connection
func (u *User) CloseConn() error {
	return u.conn.Close()
}

// IsValid return false if connCtx is nil
func (u *User) IsValid() bool {
	return u.conn != nil
}

// Receive message
func (u User) Receive(msg []byte) error {
	return u.conn.WriteBinary(msg)
}
