package core

import (
	"container/list"
	pq "github.com/lib/pq"
	"log"
	"time"
)

const (
	PubSubListenContinue = 1
	PubSubListenStop     = 0

	ProcessStatusInit   = 0
	ProcessStatusUpdate = 1
	ProcessStatusDone   = 2
	ProcessStatusError  = 3
)

type Listener interface {
	Handle(notification *pq.Notification, manager NodeManager) (int, error)
}

type SubscriberHander func(notification *pq.Notification) (int, error)

type ModelEvent struct {
	Subject  string    `json:"subject"`
	Action   string    `json:"action"`
	Type     string    `json:"type"`
	Revision int       `json:"revision"`
	Date     time.Time `json:"date"`
	Extra    string    `json:"extra"`
	Name     string    `json:"name"`
}

func NewSubscriber(conninfo string, logger *log.Logger) *Subscriber {
	return &Subscriber{
		conninfo: conninfo,
		handlers: make(map[string]*list.List, 1024),
		exit:     make(chan int),
		logger:   logger,
	}
}

type Subscriber struct {
	conninfo string
	handlers map[string]*list.List
	listener *pq.Listener
	exit     chan int
	init     bool
	logger   *log.Logger
}

func (s *Subscriber) Stop() {
	s.logger.Printf("Sending a stop to channel subscriber\n")

	s.exit <- 1
	s.listener.Close()
}

func (s *Subscriber) register() {
	if s.init {
		return
	}

	s.init = true

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			s.logger.Println(err.Error())
		}
	}

	// listen to the specific channel
	s.listener = pq.NewListener(s.conninfo, 10*time.Second, time.Minute, reportProblem)

	go s.waitAndDispatch()
}

func (s *Subscriber) waitAndDispatch() {

	// iterate over received notifications, for now, we start only one consumer with no concurrence
	for {
		select {
		case notification := <-s.listener.Notify:

			if notification == nil {
				s.logger.Printf("Subscriber: received a nil notification, the underlying driver reconnect\n", notification.Channel)
				continue
			}

			s.logger.Printf("Subscriber: received notification on channel = %s\n", notification.Channel)

			if _, ok := s.handlers[notification.Channel]; ok {
				// go some handlers register
				for e := s.handlers[notification.Channel].Front(); e != nil; e = e.Next() {

					// TODO: add a goroutine to run the handler in background
					// do something with e.Value
					if state, _ := e.Value.(SubscriberHander)(notification); state != PubSubListenContinue {
						// close listener
						s.handlers[notification.Channel].Remove(e)

						s.logger.Printf("Subscriber: removing on handler for channel = `%s` - state != PubSubListenContinue\n", notification.Channel)
					}
				}
			} else {
				s.logger.Printf("Subscriber: skipping, no handler for channel = `%s`\n", notification.Channel)
			}

		case <-time.After(20 * time.Second):
			go func() {
				s.listener.Ping()
			}()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			s.logger.Print("Subscriber: received no work for 20 seconds, checking for new work\n")

		case <-s.exit:
			return
		}

	}
}

func (s *Subscriber) ListenMessage(name string, handler SubscriberHander) {
	s.register()

	if _, ok := s.handlers[name]; !ok {
		s.handlers[name] = list.New()

		err := s.listener.Listen(name)

		if err != nil {
			panic(err)
		}
	}

	s.handlers[name].PushBack(handler)
}
