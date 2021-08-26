package eventbus

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel
type DataChannel chan DataEvent

// DataChannelSlice
type DataChannelSlice []DataChannel

// EventBus stores
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.subscribers[topic]; found {
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

var eb = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}

func PrintDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func PublisTo(topic string, data string) {
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

//func doItNow() {
//	ch1 := make(chan DataEvent)
//	ch2 := make(chan DataEvent)
//	ch3 := make(chan DataEvent)

//	eb.Subscribe("topic1", ch1)
//	eb.Subscribe("topic2", ch2)
//	eb.Subscribe("topic2", ch3)

//	go publisTo("topic1", "Hi topic 1")
//	go publisTo("topic2", `{"EleKey":"EleValue"}`)

//	for {
//		select {
//		case d := <-ch1:
//			go printDataEvent("ch1", d)
//		case d := <-ch2:
//			go printDataEvent("ch2", d)
//		case d := <-ch3:
//			go printDataEvent("ch3", d)
//		}
//	}
//}
