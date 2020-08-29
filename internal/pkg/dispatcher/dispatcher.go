package dispatcher

import (
	"context"

	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"go.uber.org/fx"
)

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

type Dispatcher struct {
	rms *Rooms
}

func (d *Dispatcher) RegisterRoom(ctx context.Context, room *entity.Room) error {
	d.rms.CreateRoom(room)
	return nil
}

func (d *Dispatcher) CloseRoom(ctx context.Context, roomUID string) error {
	room := d.rms.DeleteRoom(roomUID)
	return room.SendCloseMsg()
}

func (d *Dispatcher) JoinUserToRoom(ctx context.Context, roomUID string, user *entity.User) error {
	if !user.IsValid() {

	}

	room := d.GetRoom(ctx, roomUID)
	room.JoinUser(user)
	return nil
}

func (d *Dispatcher) DispatchMessage(ctx context.Context, msg *entity.Message) error {
	room := d.GetRoom(ctx, msg.RoomUID)
	if room == nil {
		return nil
	}
	room.PublishMessage(msg.Data)
	return nil
}

func (d *Dispatcher) GetRoom(ctx context.Context, roomUID string) *entity.Room {
	room := d.rms.GetRoom(roomUID)
	if room == nil {
		return nil
	}
	return room
}

func (d *Dispatcher) RemoveUserFromRoom(ctx context.Context, roomUID, userUID string) error {
	room := d.rms.GetRoom(roomUID)
	if room == nil {
		return nil
	}

	room.RemoveUser(userUID)
	return nil
}
