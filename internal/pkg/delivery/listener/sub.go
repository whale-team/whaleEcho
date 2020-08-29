package listener

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

// Subscriber define subscribe behavior
type Subscriber interface {
	Subscribe(subject string, msgHandler func(ctx context.Context, data []byte) error) error
	SubscribeQueue(subject, group string, msgCallback func(ctx context.Context, data []byte) error) error
}

// Listen register Listener on subjects
func Listen(ln Listener) error {
	groupName, err := os.Hostname()
	if err != nil {
		return err
	}
	err = ln.sub.Subscribe(subjects.OpenRoomSubject, ln.CreateRoom)
	log.Info().Msgf("listener: subscribe on %s", subjects.OpenRoomSubject)
	if err != nil {
		return err
	}

	err = ln.sub.Subscribe(subjects.CloseRoomSubject, ln.CloseRoom)
	log.Info().Msgf("listener: subscribe on %s", subjects.CloseRoomSubject)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		if err := ln.sub.SubscribeQueue(subjects.RoomMsgSubject, groupName, ln.DispatchMessage); err != nil {
			return err
		}
	}
	log.Info().Msgf("listener: queue subscribe on %s, queue:%s", subjects.RoomMsgSubject, groupName)
	return nil
}
