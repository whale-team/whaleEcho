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

// roomsContainer represent contaienr to contain runtime rooms
type roomsContainer struct {
	rooms map[string]*entity.Room
	mu    sync.RWMutex
	wg    *sync.WaitGroup
}

// AddRoom add new room to container, meanwhile run the room in runtime
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

// JoinRoom add participant to specific room
func (con *roomsContainer) JoinRoom(roomUID string, p entity.Participant) error {
	con.mu.RLock()
	defer con.mu.RUnlock()

	if _, ok := con.rooms[roomUID]; !ok {
		return nil
	}

	con.rooms[roomUID].Join(p)
	return nil
}

// HasRoom check if room exist
func (con *roomsContainer) HasRoom(roomUID string) bool {
	con.mu.RLock()
	defer con.mu.RUnlock()
	_, ok := con.rooms[roomUID]
	return ok
}

// Size return room size
func (con *roomsContainer) Size() int64 {
	con.mu.RLock()
	defer con.mu.RUnlock()

	return int64(len(con.rooms))
}

// RemoveRoom remove room from container, it expect this room is closed
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

// LeaveRoom remove participant from a specific room
func (con *roomsContainer) LeaveRoom(roomUID string, p entity.Participant) {
	con.mu.RLock()
	defer con.mu.RUnlock()

	if _, ok := con.rooms[roomUID]; !ok {
		return
	}

	con.rooms[roomUID].Leave(p)
}

func (con *roomsContainer) LeaveAllRooms(p entity.Participant) {
	con.mu.RLock()
	defer con.mu.RUnlock()

	for _, room := range con.rooms {
		room.Leave(p)
	}
}
