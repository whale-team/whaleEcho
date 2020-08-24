package listener

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/whale-team/whaleEcho/pkg/subjects"
)

type Subscriber interface {
	Subscribe(subject string, msgHandler func(ctx context.Context, data []byte) error) error
	SubscribeQueue(subject, group string, msgCallback func(ctx context.Context, data []byte) error) error
}

func Listen(ln Listener) error {
	groupName, err := os.Hostname()
	if err != nil {
		return err
	}
	err = ln.sub.Subscribe(subjects.OpenRoomSubject, ln.CreateRoom)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		if err := ln.sub.SubscribeQueue(subjects.RoomMsgSubject, groupName, ln.DispatchMessage); err != nil {
			return err
		}
	}
	log.Info().Msgf("listener: subscribe %s", subjects.RoomMsgSubject)
	return nil
}
