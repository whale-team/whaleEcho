package entity

import (
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/internal/pkg/app/entity/value"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

const (
	maxWorker = 5
)

// NewRoom construct room struct
func NewRoom() *Room {
	return &Room{
		Participants: sync.Map{},
		closeSignal:  make(chan struct{}),
		workers:      workerpool.New(30),
		workerSize:   30,
		wg:           &sync.WaitGroup{},
		sysMessage:   sysMessageSet,
	}
}

// Participant define interface allow room to communicate with
type Participant interface {
	Receive(msg MsgData) error
	GetID() int64
}

// Room represent chating room
type Room struct {
	ID          int64
	UID         string
	Limit       int64
	CreatorID   int64
	CreatorName string

	Participants sync.Map
	Subscribe    Subscriber
	msgCh        <-chan *nats.Msg
	joinCh       <-chan *nats.Msg
	closeSignal  chan struct{}
	workers      *workerpool.WorkerPool
	workerSize   int64
	mu           sync.RWMutex
	closed       bool
	sysMessage   map[value.SysMsgType]*SysMessage
	wg           *sync.WaitGroup
}

// Subject to room pub sub subject
func (r *Room) Subject() string {
	return subjects.RoomSubject(r.UID)
}

// Join add paricipant to this room
func (r *Room) Join(p Participant) {
	r.Participants.Store(p.GetID(), p)
}

// Leave remove participant from room
func (r *Room) Leave(p Participant) {
	r.Participants.Delete(p.GetID())
}

// SetMsgChannel set msg channel, this channel receive message from nats connection
func (r *Room) SetMsgChannel(ch <-chan *nats.Msg) {
	r.msgCh = ch
}

// Run start room
func (r *Room) Run() {
	go r.run()
}

// PushMessage push message to paricipats
func (r *Room) PushMessage(msg *Message) {
	r.Participants.Range(func(key, val interface{}) bool {
		r.workers.Submit(func() {
			receiver := val.(Participant)
			if err := receiver.Receive(msg); err != nil {
				log.Error().Err(err).Msgf("room: PushMessage to participant failed, p_id:%d", receiver.GetID())
			}
		})
		return true
	})
}

// Close room
func (r *Room) Close() {
	close(r.closeSignal)
	r.Subscribe.Unsubscribe()
	r.wg.Wait()
}

func (r *Room) run() {
	r.wg.Add(1)
	defer r.wg.Done()
	for {
		select {
		case msg := <-r.msgCh:
			r.PushMessage(&Message{Msg: msg})
		case <-r.closeSignal:
			r.closed = true
			r.notifyClose()
			return
		}
	}
}

func (r *Room) notifyClose() {
	r.Participants.Range(func(key, val interface{}) bool {
		receiver := val.(Participant)
		if err := receiver.Receive(r.sysMessage[value.CloseRoom]); err != nil {
			log.Error().Err(err).Msgf("room: Notify Close Msg to participant failed, p_id:%d", receiver.GetID())
		}
		return true
	})
}

// Subscriber ...
type Subscriber interface {
	Unsubscribe() error
	Drain() error
	IsValid() bool
}
