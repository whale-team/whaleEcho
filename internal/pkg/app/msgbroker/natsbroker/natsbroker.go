package natsbroker

import (
	"github.com/whale-team/whaleEcho/pkg/natspool"
)

// NatsBroker ...
type NatsBroker struct {
	Client natspool.Client
}

// // Publish ...
// func (broker NatsBroker) Publish(ctx context.Context, topic model.Topic, msg model.Message) error {
// 	conn, err := broker.Client.Conn()
// 	if err != nil {
// 		return errors.Wrap(err, "broker#Publish: get client connection failed")
// 	}
// 	data, err := json.Marshal(&msg)
// 	if err != nil {
// 		return errors.Wrap(err, "broker#Publish: unmarshaling message struct failed")
// 	}
// 	err = conn.Publish(topic.String(), data)
// 	if err != nil {
// 		return errors.Wrap(err, "broker#Publish: message failed")
// 	}
// 	err = broker.Client.FlushAndClose(conn, false)
// 	if err != nil {
// 		return errors.Wrap(err, "broker#Publish: publish message failed")
// 	}
// 	if err := conn.LastError(); err != nil {
// 		return errors.Wrap(err, "broker:Publish: conn encounter unknown filed")
// 	}
// 	return nil
// }

// // Subscribe ...
// func (broker NatsBroker) Subscribe(ctx context.Context, topic model.Topic, callback app.SubCallback) (app.ConnCloser, error) {
// 	conn, err := broker.Client.Conn()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "broker#Subsribe: get client connection failed")
// 	}
// 	_, err = conn.Subscribe(topic.String(), func(message *nats.Msg) {
// 		msg := model.Message{}
// 		err := json.Unmarshal(message.Data, &msg)
// 		callback(ctx, msg, err)
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "broker#Subsribe: subscribe topic failed")
// 	}

// 	if err := conn.LastError(); err != nil {
// 		return nil, errors.Wrap(err, "broker#Subsribe: conn encounter unknown failed")
// 	}

// 	return conn, nil
// }
