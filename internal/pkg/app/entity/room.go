package entity

import (
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/nats-io/nats.go"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

const (
	maxWorker = 5
)

func NewRoom() *Room {
	return &Room{
		Participants: sync.Map{},
		closeSignal:  make(chan struct{}),
		workers:      workerpool.New(30),
		workerSize:   30,
		wg:           &sync.WaitGroup{},
		closedMsg:    RoomCloseMessage,
	}
}

type Participant interface {
	Receive(msg *Message) error
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
	closedMsg    *Message
	wg           *sync.WaitGroup
}

func (r *Room) Subject() string {
	return subjects.RoomSubject(r.UID)
}

func (r *Room) SetClosedMsg(msg *Message) {
	r.closedMsg = msg
}

func (r *Room) Join(p Participant) {
	r.Participants.Store(p.GetID(), p)
}

func (r *Room) Leave(p Participant) {
	r.Participants.Delete(p.GetID())
}

func (r *Room) SetMsgChannel(ch <-chan *nats.Msg) {
	r.msgCh = ch
}

func (r *Room) Run() {
	go r.run()
}

func (r *Room) PushMessage(msg *Message) {
	r.Participants.Range(func(key, val interface{}) bool {
		r.workers.Submit(func() {
			receiver := val.(Participant)
			receiver.Receive(msg)
		})
		return true
	})
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
		receiver.Receive(r.closedMsg)
		return true
	})
}

func (r *Room) Close() {
	close(r.closeSignal)
	r.Subscribe.Unsubscribe()
	r.wg.Wait()
}

// Subscriber ...
type Subscriber interface {
	Unsubscribe() error
	Drain() error
	IsValid() bool
}
