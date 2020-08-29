package dispatcher

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"go.uber.org/fx"
)

// Params Dispatcher dependencies
type Params struct {
	fx.In
	Rooms *Rooms
}

// New create singleton dispatcher
func New(params Params) app.Dispatcher {
	return &Dispatcher{
		rms: params.Rooms,
	}
}

// Dispatcher dispatch message to client connections
type Dispatcher struct {
	rms *Rooms
}

// RegisterRoom register a room client container
func (d *Dispatcher) RegisterRoom(ctx context.Context, room *entity.Room) error {
	d.rms.CreateRoom(room)
	return nil
}

// CloseRoom delete the room from memory and send close message to connection
func (d *Dispatcher) CloseRoom(ctx context.Context, roomUID string) error {
	room := d.rms.DeleteRoom(roomUID)
	if err := room.SendCloseMsg(); err != nil {
		return err
	}
	room = nil
	return nil
}

// JoinUserToRoom join a client connection to room
func (d *Dispatcher) JoinUserToRoom(ctx context.Context, roomUID string, user *entity.User) error {
	if !user.IsValid() {

	}

	room := d.GetRoom(ctx, roomUID)
	room.JoinUser(user)
	return nil
}

// DispatchMessage dispatch message to all client connections
func (d *Dispatcher) DispatchMessage(ctx context.Context, msg *entity.Message) error {
	room := d.GetRoom(ctx, msg.RoomUID)
	if room == nil {
		return nil
	}
	room.PublishMessage(msg.Data)
	return nil
}

// GetRoom get the room from memory
func (d *Dispatcher) GetRoom(ctx context.Context, roomUID string) *entity.Room {
	room := d.rms.GetRoom(roomUID)
	if room == nil {
		return nil
	}
	return room
}

// RemoveUserFromRoom remove client connection
func (d *Dispatcher) RemoveUserFromRoom(ctx context.Context, roomUID, userUID string) error {
	room := d.rms.GetRoom(roomUID)
	if room == nil {
		return nil
	}

	room.RemoveUser(userUID)
	return nil
}
