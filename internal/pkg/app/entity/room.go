package entity

import (
	"sync"
)

const (
	maxWorker = 5
)

// NewRoom construct room struct
func NewRoom() *Room {
	return &Room{}
}

// Room represent chating room
type Room struct {
	UID          string
	MembersLimit int64
	Name         string
	MembersCount int64

	usersMap map[string]*User
	rw       sync.RWMutex
}

func (room *Room) getUsersMap() map[string]*User {
	if room.usersMap == nil {
		room.usersMap = make(map[string]*User)
	}
	return room.usersMap
}

func (room *Room) CurrentMembersCount() int {
	room.rw.RLock()
	defer room.rw.RUnlock()
	return len(room.getUsersMap())
}

func (room *Room) JoinUser(user *User) {
	room.rw.Lock()
	defer room.rw.Unlock()

	room.getUsersMap()[user.UID] = user
}

func (room *Room) RemoveUser(userUID string) *User {
	room.rw.Lock()
	defer room.rw.Unlock()

	user := room.getUsersMap()[userUID]
	delete(room.usersMap, userUID)
	return user
}

func (room *Room) PublishMessage(msg []byte) error {
	room.rw.Lock()
	defer room.rw.Unlock()

	for _, u := range room.getUsersMap() {
		u.Receive(msg)
	}

	return nil
}

func (room *Room) SendCloseMsg() error {
	closeRoomMsg := []byte("Room" + room.Name + "Closed")
	room.rw.Lock()
	defer room.rw.Unlock()

	for _, u := range room.getUsersMap() {
		u.Receive(closeRoomMsg)
	}
	return nil
}
