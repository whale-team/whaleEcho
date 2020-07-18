package roomscenter

import (
	"sync"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

func newContainer() *roomsContainer {
	return &roomsContainer{
		rooms: make(map[string]*entity.Room),
		wg:    &sync.WaitGroup{},
	}
}

type roomsContainer struct {
	rooms map[string]*entity.Room
	mu    sync.RWMutex
	wg    *sync.WaitGroup
}

func (con *roomsContainer) AddRoom(room *entity.Room) error {
	con.mu.Lock()
	defer con.mu.Unlock()

	if _, ok := con.rooms[room.UID]; ok {
		return nil
	}

	con.rooms[room.UID] = room
	con.rooms[room.UID].Run()
	return nil
}

func (con *roomsContainer) JoinRoom(roomUID string, p entity.Participant) error {
	con.mu.RLock()
	defer con.mu.RUnlock()

	if _, ok := con.rooms[roomUID]; !ok {
		return nil
	}

	con.rooms[roomUID].Join(p)
	return nil
}

func (con *roomsContainer) HasRoom(roomUID string) bool {
	con.mu.RLock()
	defer con.mu.RUnlock()
	_, ok := con.rooms[roomUID]
	return ok
}

func (con *roomsContainer) Size() int64 {
	con.mu.RLock()
	defer con.mu.RUnlock()

	return int64(len(con.rooms))
}

func (con *roomsContainer) RemoveRoom(roomUID string) error {
	con.mu.Lock()
	defer con.mu.Unlock()

	if _, ok := con.rooms[roomUID]; !ok {
		return nil
	}

	room := con.rooms[roomUID]
	delete(con.rooms, roomUID)
	room.Close()
	return nil
}

func (con *roomsContainer) LeaveRoom(roomUID string, p entity.Participant) {
	con.mu.RLock()
	defer con.mu.RUnlock()

	if _, ok := con.rooms[roomUID]; !ok {
		return
	}

	con.rooms[roomUID].Leave(p)
}
