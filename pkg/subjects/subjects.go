package subjects

const (
	OpenRoomSubject  = "rooms.open"
	CloseRoomSubject = "rooms.close"
)

func RoomSubject(roomUID string) string {
	return "room." + roomUID
}
