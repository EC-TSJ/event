package event

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type imessage interface {
	// The message to be sent to clients
	prepareMessage() []byte
}

type teventMessage struct {
	id    string
	event string
	data  string
}

func (m *teventMessage) prepareMessage() []byte {
	var data bytes.Buffer
	if len(m.id) > 0 {
		data.WriteString(fmt.Sprintf("id: %s\n", strings.Replace(m.id, "\n", "", -1)))
	}
	if len(m.event) > 0 {
		data.WriteString(fmt.Sprintf("event: %s\n", strings.Replace(m.event, "\n", "", -1)))
	}
	if len(m.data) > 0 {
		lines := strings.Split(m.data, "\n")
		for _, line := range lines {
			data.WriteString(fmt.Sprintf("data: %s\n", line))
		}
	}
	data.WriteString("\n")
	return data.Bytes()
}

type retryMessage struct {
	retry time.Duration
}

func (m *retryMessage) prepareMessage() []byte {
	return []byte(fmt.Sprintf("retry: %d\n\n", m.retry/time.Millisecond))
}

// EventSource interface provides methods for sending messages and closing all connections.
type IEventSource interface {
	// it should implement ServerHTTP method
	http.Handler

	// send message to all consumers
	SendEventMessage(data, event, id string)

	// send retry message to all consumers
	SendRetryMessage(duration time.Duration)

	// consumers count
	Users() int

	// close and clear all consumers
	Close()
}

type teventSource struct {
	customHeadersFunc func(*http.Request) [][]byte
	sink              chan imessage
	staled            chan *tuser
	add               chan *tuser
	close             chan bool
	idleTimeout       time.Duration
	retry             time.Duration
	timeout           time.Duration
	closeOnTimeout    bool
	gzip              bool
	users             *list.List
}

func (es *teventSource) Close() {
	es.close <- true
}

// ServeHTTP implements http.Handler interface.
func (es *teventSource) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	cons, err := new_user(resp, req, es)
	if err != nil {
		log.Print("No puedo crear conexión a un usuario: ", err)
		return
	}
	es.add <- cons
}

func (es *teventSource) sendMessage(m imessage) {
	es.sink <- m
}

func (es *teventSource) SendEventMessage(data, event, id string) {
	em := &teventMessage{id, event, data}
	es.sendMessage(em)
}

func (es *teventSource) SendRetryMessage(t time.Duration) {
	es.sendMessage(&retryMessage{t})
}

func (es *teventSource) Users() int {
	return es.users.Len()
}

type TSettings struct {
	// SetTimeout sets the write timeout for individual messages. The
	// default is 2 seconds.
	Timeout time.Duration

	// CloseOnTimeout sets whether a write timeout should close the
	// connection or just drop the message.
	//
	// If the connection gets closed on a timeout, it's the client's
	// responsibility to re-establish a connection. If the connection
	// doesn't get closed, messages might get sent to a potentially dead
	// client.
	//
	// The default is true.
	CloseOnTimeout bool

	// Sets the timeout for an idle connection. The default is 30 minutes.
	IdleTimeout time.Duration

	// Gzip sets whether to use gzip Content-Encoding for clients which
	// support it.
	//
	// The default is false.
	Gzip bool
}

func DefaultSettings() *TSettings {
	return &TSettings{
		Timeout:        2 * time.Second,
		CloseOnTimeout: true,
		IdleTimeout:    30 * time.Minute,
		Gzip:           false,
	}
}

func controlProcess(es *teventSource) {
	for {
		select {
		case em := <-es.sink:
			message := em.prepareMessage()
			func() {
				for e := es.users.Front(); e != nil; e = e.Next() {
					c := e.Value.(*tuser)

					// Only send this message if the consumer isn't staled
					if !c.staled {
						select {
						case c.in <- message:
						default:
						}
					}
				}
			}()
		case <-es.close:
			close(es.sink)
			close(es.add)
			close(es.staled)
			close(es.close)

			func() {
				for e := es.users.Front(); e != nil; e = e.Next() {
					c := e.Value.(*tuser)
					close(c.in)
				}
			}()

			es.users.Init()
			return
		case c := <-es.add:
			func() {
				es.users.PushBack(c)
			}()
		case c := <-es.staled:
			toRemoveEls := make([]*list.Element, 0, 1)
			func() {
				for e := es.users.Front(); e != nil; e = e.Next() {
					if e.Value.(*tuser) == c {
						toRemoveEls = append(toRemoveEls, e)
					}
				}
			}()
			func() {
				for _, e := range toRemoveEls {
					es.users.Remove(e)
				}
			}()
			close(c.in)
		}
	}
}

// New creates new EventSource instance.
func EventSource(settings *TSettings, customHeadersFunc func(*http.Request) [][]byte) IEventSource {
	if settings == nil {
		settings = DefaultSettings()
	}

	es := new(teventSource)
	es.customHeadersFunc = customHeadersFunc
	es.sink = make(chan imessage, 1)
	es.close = make(chan bool)
	es.staled = make(chan *tuser, 1)
	es.add = make(chan *tuser)
	es.users = list.New()
	es.timeout = settings.Timeout
	es.idleTimeout = settings.IdleTimeout
	es.closeOnTimeout = settings.CloseOnTimeout
	es.gzip = settings.Gzip
	go controlProcess(es)

	return es
}
