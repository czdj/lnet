package lutils

//type Event interface{
//	GetEventId() int32
//	GetEventMsg() interface{}
//}

type Event struct {
	Id  int32
	Msg interface{}
}

func (this *Event) GetEventId() int32 {
	return this.Id
}
func (this *Event) GetEventMsg() interface{} {
	return this.Msg
}

//订阅者
type ISubscriber interface {
	Subscribe(pubsub *PubSub, realSub ISubscriber, eventId int32)
	Unsubscribe(pubsub *PubSub, realSub ISubscriber, eventId int32)
	HandleMsg(event *Event)
}

type Subscriber struct {
}

func (this *Subscriber) Subscribe(pubsub *PubSub, realSub ISubscriber, eventId int32) {
	pubsub.AddSubscriber(realSub, eventId)
}

func (this *Subscriber) Unsubscribe(pubsub *PubSub, realSub ISubscriber, eventId int32) {
	pubsub.RemoveSubscriber(realSub, eventId)
}

func (this *Subscriber) HandleMsg(event *Event) {

}

//发布者
type IPublisher interface {
	Publish(pubsub *PubSub, event *Event)
}

type Publisher struct {
}

func (this *Publisher) Publish(pubsub *PubSub, event *Event) {
	pubsub.Publish(event)
}

//发布订阅管理类
type PubSub struct {
	Event2SubscriberMap map[int32][]ISubscriber
}

var g_pubSubMag *PubSub

func GetPubSubMng() *PubSub {
	return g_pubSubMag
}
func NewPubSub() *PubSub {
	re := &PubSub{
		Event2SubscriberMap: make(map[int32][]ISubscriber),
	}

	return re
}

func (this *PubSub) AddSubscriber(sub ISubscriber, eventId int32) {
	v, ok := this.Event2SubscriberMap[eventId]
	if !ok {
		this.Event2SubscriberMap[eventId] = []ISubscriber{sub}
		return
	}

	for _, vv := range v {
		if vv == sub {
			return
		}
	}

	v = append(v, sub)
}

func (this *PubSub) RemoveSubscriber(sub ISubscriber, eventId int32) {
	v, ok := this.Event2SubscriberMap[eventId]
	if !ok {
		return
	}

	for k, vv := range v {
		if vv == sub {
			v = append(v[:k], v[k+1:]...)
			break
		}
	}

	if len(v) == 0 {
		delete(this.Event2SubscriberMap, eventId)
	}
}

func (this *PubSub) Publish(event *Event) {
	v, ok := this.Event2SubscriberMap[event.Id]
	if !ok {
		return
	}

	for _, vv := range v {
		vv.HandleMsg(event)
	}
}

func init() {
	g_pubSubMag = NewPubSub()
}
