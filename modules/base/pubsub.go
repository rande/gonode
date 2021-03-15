// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	pq "github.com/lib/pq"
	"github.com/rande/gonode/core/helper"
	log "github.com/sirupsen/logrus"
)

const (
	PubSubListenContinue = 1
	PubSubListenStop     = 0

	ProcessStatusInit   = 0  // default value, nothing to do
	ProcessStatusReady  = 1  // process is ready to be handled
	ProcessStatusUpdate = 2  // update in progress
	ProcessStatusDone   = 3  // done, can also be set to init. Done also mean the related task cannot be restarted
	ProcessStatusError  = -1 // an error occurs
)

type Listener interface {
	Handle(notification *pq.Notification, manager NodeManager) (int, error)
}

type SubscriberHander func(notification *pq.Notification) (int, error)

type ModelEvent struct {
	Subject     string    `json:"subject"`
	Action      string    `json:"action"`
	Type        string    `json:"type"`
	Revision    int       `json:"revision"`
	Date        time.Time `json:"date"`
	Extra       string    `json:"extra"`
	Name        string    `json:"name"`
	NewRevision bool      `json:"new_revision"`
}

func NewSubscriber(conninfo string, logger *log.Logger) *Subscriber {
	return &Subscriber{
		conninfo: conninfo,
		handlers: make(map[string]*list.List, 1024),
		exit:     make(chan int),
		logger:   logger,
		channels: make([]string, 0),
	}
}

func CreateModelEvent(notification *pq.Notification) *ModelEvent {
	m := &ModelEvent{}

	json.Unmarshal([]byte(notification.Extra), m)

	return m
}

type Subscriber struct {
	conninfo string
	handlers map[string]*list.List
	listener *pq.Listener
	exit     chan int
	init     bool
	logger   *log.Logger
	channels []string
}

func (s *Subscriber) Stop() {
	s.logger.WithFields(log.Fields{
		"module": "node.pubsub",
	}).Debug("Sending a stop to channel subscriber")

	s.exit <- 1
	s.listener.Close()
}

func (s *Subscriber) Register() {
	if s.init {
		return
	}

	s.init = true

	// listen to the specific channel
	s.listener = pq.NewListener(s.conninfo, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			s.logger.Println(err.Error())
		}
	})

	for _, name := range s.channels {
		err := s.listener.Listen(name)
		helper.PanicOnError(err)
	}

	go s.waitAndDispatch()
}

func (s *Subscriber) waitAndDispatch() {
	// iterate over received notifications, for now, we start only one consumer with no concurrence
	for {
		select {
		case notification := <-s.listener.Notify:

			if notification == nil {
				s.logger.WithFields(log.Fields{
					"module": "node.pubsub",
				}).Warn("received a nil notification, the underlying driver reconnect")

				continue
			}

			s.logger.WithFields(log.Fields{
				"channel": notification.Channel,
				"module":  "node.pubsub",
			}).Debug("received notification on channel")

			if _, ok := s.handlers[notification.Channel]; ok {
				// go some handlers register
				for e := s.handlers[notification.Channel].Front(); e != nil; e = e.Next() {
					go func(e *list.Element) {
						var f = e.Value.(SubscriberHander)
						s.logger.WithFields(log.Fields{
							"channel": notification.Channel,
							"payload": notification.Extra,
							"module":  "node.pubsub",
							"handler": fmt.Sprintf("%T", f),
						}).Debug("send payload to handler")

						if state, err := f(notification); state != PubSubListenContinue {
							// close listener
							s.handlers[notification.Channel].Remove(e)
							s.logger.WithFields(log.Fields{
								"channel": notification.Channel,
								"state":   state,
								"module":  "node.pubsub",
							}).Debug("removing handler for channel - state != PubSubListenContinue")
						} else if err != nil {
							s.logger.WithFields(log.Fields{
								"channel": notification.Channel,
								"payload": notification.Extra,
								"module":  "node.pubsub",
								"error":   err.Error(),
							}).Debug("End processing message (ie: func return, goroutine started ?)")
						}

						s.logger.WithFields(log.Fields{
							"channel": notification.Channel,
							"payload": notification.Extra,
							"module":  "node.pubsub",
							"handler": fmt.Sprintf("%T", f),
						}).Debug("End processing message (ie: func return, goroutine started ?)")

					}(e)
				}
			} else {
				s.logger.WithFields(log.Fields{
					"channel": notification.Channel,
					"module":  "node.pubsub",
				}).Debug("skipping, no handler for channel")
			}

		case <-time.After(20 * time.Second):
			go func() {
				s.listener.Ping()
			}()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			s.logger.WithFields(log.Fields{
				"module": "node.pubsub",
			}).Debug("received no work for 20 seconds, checking for new work")

		case <-s.exit:
			return
		}
	}
}

func (s *Subscriber) ListenMessage(name string, handler SubscriberHander) {
	if _, ok := s.handlers[name]; !ok {
		s.handlers[name] = list.New()

		if s.init {
			err := s.listener.Listen(name)
			helper.PanicOnError(err)
		} else {
			s.channels = append(s.channels, name)
		}
	}

	s.handlers[name].PushBack(handler)
}
