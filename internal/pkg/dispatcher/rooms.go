package dispatcher

import (
	"sync"

	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
)

var _rooms *Rooms

var once = sync.Once{}

// NewRooms build a rooms connection container
func NewRooms() *Rooms {
	once.Do(func() {
		_rooms = &Rooms{
			roomsMap: make(map[string]*entity.Room),
		}
	})

	return _rooms
}

// Rooms is runtime rooms info storer to keep alive connection between room and user
type Rooms struct {
	roomsMap map[string]*entity.Room
	mu       sync.RWMutex
}

// Clear clear all rooms
func (c *Rooms) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, r := range c.roomsMap {
		r.Clear()
	}

	c.Reset()
}

// Len number of rooms
func (c *Rooms) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.roomsMap)
}

// Reset reset all rooms
func (c *Rooms) Reset() {
	c.roomsMap = make(map[string]*entity.Room)
}

// GetRoom get room based on roomUID
func (c *Rooms) GetRoom(roomUID string) *entity.Room {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.roomsMap[roomUID]
}

// CreateRoom create room into memory
func (c *Rooms) CreateRoom(room *entity.Room) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.roomsMap[room.UID]; !ok {
		c.roomsMap[room.UID] = room
	}
}

// DeleteRoom delete room from memory
func (c *Rooms) DeleteRoom(roomUID string) *entity.Room {
	c.mu.Lock()
	defer c.mu.Unlock()

	room := c.roomsMap[roomUID]
	delete(c.roomsMap, roomUID)
	return room
}
