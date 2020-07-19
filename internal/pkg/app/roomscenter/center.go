package roomscenter

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity"
	"github.com/whale-team/whaleEcho/internal/pkg/app/msgbroker"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

var once = sync.Once{}

var center *Center

// Center represent rooms center to manage all runtime rooms
type Center struct {
	*roomsContainer
	openCh       chan *nats.Msg
	closeCh      chan *nats.Msg
	createSub    entity.Subscriber
	deleteSub    entity.Subscriber
	broker       msgbroker.MsgBroker
	mu           *sync.RWMutex
	ctx          context.Context
	ctxCancel    context.CancelFunc
	asyncHandler AsyncHandler
	binded       bool
}

// New build a rooms center to manage all rooms
func New(broker msgbroker.MsgBroker) (*Center, error) {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		center = &Center{
			roomsContainer: newContainer(),
			openCh:         make(chan *nats.Msg, 1),
			closeCh:        make(chan *nats.Msg, 1),
			ctx:            ctx,
			ctxCancel:      cancel,
			asyncHandler:   NewDefaultHandler(broker),
		}
	})
	if err := center.bindSubscribe(broker); err != nil {
		return nil, err
	}
	center.Start()
	return center, nil
}

// SetAsyncHandler setup an async handler
func (center *Center) SetAsyncHandler(handler AsyncHandler) {
	center.asyncHandler = handler
}

// Start start rooms center
//  Rooms Center will listen open room and close room event from nats
func (center *Center) Start() {
	go center.listen()
}

// Shutdown shutdown the rooms
func (center *Center) Shutdown() {
	center.ctxCancel()
}

func (center *Center) listen() {
	for {
		select {
		case msg := <-center.openCh:
			room := entity.NewRoom()
			if err := center.asyncHandler.OpenRoom(center.ctx, msg, room); err != nil {
				center.asyncHandler.ErrHandle(err, msg, room)
			}
			center.AddRoom(room)
		case msg := <-center.closeCh:
			room := &entity.Room{}
			err := center.asyncHandler.CloseRoom(center.ctx, msg, room)
			if err != nil {
				center.asyncHandler.ErrHandle(err, msg, room)
				continue
			}
			center.RemoveRoom(room.UID)
		case <-center.ctx.Done():
			break
		}
	}
}

func (center *Center) bindSubscribe(broker msgbroker.MsgBroker) error {
	if center.binded {
		return nil
	}

	sub, err := broker.BindChannelWithSubject(center.ctx, subjects.OpenRoomSubject, center.openCh)
	if err != nil {
		return err
	}
	center.createSub = sub
	sub, err = broker.BindChannelWithSubject(center.ctx, subjects.CloseRoomSubject, center.closeCh)
	if err != nil {
		return err
	}
	center.deleteSub = sub
	center.binded = true
	return nil
}
