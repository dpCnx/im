package rabbitmq

import (
	"github.com/google/wire"
	"github.com/streadway/amqp"
	"im/api/pb"
)

var ProviderSet = wire.NewSet(NewRabbitmq, NewMsgRabbitmq)

type Rabbitmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitmq(config *pb.Data) (*Rabbitmq, error) {
	conn, err := amqp.Dial(config.Rabbitmq.Address)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Rabbitmq{conn: conn, ch: ch}, nil
}
func (mq *Rabbitmq) Close() {
	if mq.conn != nil {
		_ = mq.conn.Close()
	}
	if mq.ch != nil {
		_ = mq.ch.Close()
	}
}

type MsgRabbitmq struct {
	rmq          *Rabbitmq
	exchangeName string
	queueName    string
	key          string
}

func NewMsgRabbitmq(mq *Rabbitmq, config *pb.Data) *MsgRabbitmq {
	return &MsgRabbitmq{
		rmq:          mq,
		exchangeName: config.Rabbitmq.ExchangeName,
		queueName:    config.Rabbitmq.QueueName,
		key:          config.Rabbitmq.Key,
	}
}

func (msgRabbitmq *MsgRabbitmq) CreateExchange() error {
	return msgRabbitmq.rmq.ch.ExchangeDeclare(msgRabbitmq.exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
}

func (msgRabbitmq *MsgRabbitmq) CreateQueue() (amqp.Queue, error) {
	return msgRabbitmq.rmq.ch.QueueDeclare(msgRabbitmq.queueName, true, false, false, false, nil)
}

func (msgRabbitmq *MsgRabbitmq) BindQueue() error {
	return msgRabbitmq.rmq.ch.QueueBind(msgRabbitmq.queueName, msgRabbitmq.key, msgRabbitmq.exchangeName, false, nil)
}

func (msgRabbitmq *MsgRabbitmq) Confirm(f func(deliveryTag uint64)) error {
	go func() {
		if err := msgRabbitmq.rmq.ch.Confirm(false); err != nil {
			return
		}
		confirmChan := msgRabbitmq.rmq.ch.NotifyPublish(make(chan amqp.Confirmation))
		for cc := range confirmChan {
			if !cc.Ack {
				f(cc.DeliveryTag)
			}
		}

	}()
	return nil
}

func (msgRabbitmq *MsgRabbitmq) Return(f func(res []byte)) error {
	go func() {
		cr := msgRabbitmq.rmq.ch.NotifyReturn(make(chan amqp.Return))
		for r := range cr {
			f(r.Body)
		}
	}()
	return nil
}

func (msgRabbitmq *MsgRabbitmq) PushMsg(msg []byte) error {
	return msgRabbitmq.rmq.ch.Publish(msgRabbitmq.exchangeName, msgRabbitmq.key, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         msg,
	})
}

func (msgRabbitmq *MsgRabbitmq) Consume(f func([]byte) error) error {
	msgs, err := msgRabbitmq.rmq.ch.Consume(msgRabbitmq.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			if err = f(msg.Body); err != nil {
				_ = msg.Reject(true)
			} else {
				// multiple 批量确认
				_ = msg.Ack(false)
			}
		}
	}()
	return nil
}
