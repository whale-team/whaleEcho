package entity

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

var (
	CloseRoomMsg = &echoproto.Message{
		SenderName: "whale",
		Type:       echoproto.MessageType_Text,
		Text:       "room closed!",
	}
)

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

func (room *Room) Clear() {
	room.rw.Lock()
	defer room.rw.Unlock()

	var err error

	for _, u := range room.getUsersMap() {
		if err = u.CloseConn(); err != nil {
			log.Error().Err(err).Msgf("room: Clear close user connection occur error, err:%+v, userUID:%s", err, u.UID)
		}
	}
}

func (room *Room) SendCloseMsg() error {
	room.rw.Lock()
	defer room.rw.Unlock()

	for _, u := range room.getUsersMap() {
		u.Receive(room.getCloseRoomMsg())
	}
	return nil
}

func (room *Room) getCloseRoomMsg() []byte {
	msg := CloseRoomMsg
	msg.SentAt = time.Now().Unix()
	msg.RoomUid = room.UID
	data, _ := proto.Marshal(msg)
	return data
}
