package entity

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

var (
	// CloseRoomMsg close room system msg
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
	uids     []string
	rw       sync.RWMutex
}

func (room *Room) getUsersMap() map[string]*User {
	if room.usersMap == nil {
		room.usersMap = make(map[string]*User)
	}
	return room.usersMap
}

func (room *Room) getUids() []string {
	if room.uids == nil {
		room.uids = make([]string, 0, 1)
	}
	return room.uids
}

// CurrentMembersCount return current numner of client connections
func (room *Room) CurrentMembersCount() int {
	room.rw.RLock()
	defer room.rw.RUnlock()
	return len(room.getUsersMap())
}

// JoinUser join user to the room
func (room *Room) JoinUser(user *User) {
	room.rw.Lock()
	defer room.rw.Unlock()

	room.getUsersMap()[user.UID] = user
	room.uids = append(room.getUids(), user.UID)
}

func (room *Room) RemoveUser(userUID string) *User {
	room.rw.Lock()
	defer room.rw.Unlock()

	user := room.getUsersMap()[userUID]
	delete(room.usersMap, userUID)
	return user
}

func (room *Room) PublishMessage(msg []byte) error {
	return room.publishMessageFaster(msg)
}

func (room *Room) publishMessage(msg []byte) error {
	room.rw.Lock()
	defer room.rw.Unlock()

	for _, u := range room.getUsersMap() {
		u.Receive(msg)
	}

	return nil
}

func (room *Room) publishMessageFaster(msg []byte) error {
	room.rw.Lock()
	defer room.rw.Unlock()

	usersMap := room.getUsersMap()
	uids := room.getUids()

	uidss := DivideSlice(uids, 4)

	wg := sync.WaitGroup{}

	wg.Add(len(uidss))

	for _, uids := range uidss {
		go func(uids []string) {
			for _, uid := range uids {
				usersMap[uid].Receive(msg)
			}
			wg.Done()
		}(uids)
	}

	wg.Wait()
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

func DivideSlice(slice []string, base int) [][]string {
	l := len(slice)
	if base > l {
		return [][]string{slice}
	}

	diff := l / base
	if l%base != 0 {
		diff++
	}

	res := make([][]string, base)
	upbound, lowbound := 0, diff
	for i := range res {
		res[i] = slice[upbound:lowbound]
		upbound = lowbound
		lowbound = upbound + diff
		if lowbound >= len(slice) {
			lowbound = len(slice)
		}
	}

	return res
}
