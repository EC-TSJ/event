package event

import (
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	XAccel       = "X-Accel-Buffering: no"
	AccessOrigin = "Access-Control-Allow-Origin: *"
)

type tuser struct {
	conn   io.WriteCloser
	es     *teventSource
	in     chan []byte
	staled bool
}

type tgzipConn struct {
	net.Conn
	*gzip.Writer
}

func (gc tgzipConn) Write(b []byte) (int, error) {
	n, err := gc.Writer.Write(b)
	if err != nil {
		return n, err
	}

	return n, gc.Writer.Flush()
}

func (gc tgzipConn) Close() error {
	err := gc.Writer.Close()
	if err != nil {
		return err
	}

	return gc.Conn.Close()
}

// newConsumer
func new_user(resp http.ResponseWriter, req *http.Request, es *teventSource) (*tuser, error) {
	conn, _, err := resp.(http.Hijacker).Hijack()
	if err != nil {
		return nil, err
	}

	consumer := &tuser{
		conn:   conn,
		es:     es,
		in:     make(chan []byte, 10),
		staled: false,
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/event-stream\r\n"))
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = conn.Write([]byte("Vary: Accept-Encoding\r\n"))
	if err != nil {
		conn.Close()
		return nil, err
	}

	if es.gzip && (req == nil || strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")) {
		_, err = conn.Write([]byte("Content-Encoding: gzip\r\n"))
		if err != nil {
			conn.Close()
			return nil, err
		}

		consumer.conn = tgzipConn{conn, gzip.NewWriter(conn)}
	}

	if es.customHeadersFunc != nil {
		for _, header := range es.customHeadersFunc(req) {
			_, err = conn.Write(header)
			if err != nil {
				conn.Close()
				return nil, err
			}
			_, err = conn.Write([]byte("\r\n"))
			if err != nil {
				conn.Close()
				return nil, err
			}
		}
	}

	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		conn.Close()
		return nil, err
	}

	go func() {
		idleTimer := time.NewTimer(es.idleTimeout)
		defer idleTimer.Stop()
		for {
			select {
			case message, open := <-consumer.in:
				if !open {
					consumer.conn.Close()
					return
				}
				conn.SetWriteDeadline(time.Now().Add(consumer.es.timeout))
				_, err := consumer.conn.Write(message)
				if err != nil {
					netErr, ok := err.(net.Error)
					if !ok || !netErr.Timeout() || consumer.es.closeOnTimeout {
						consumer.staled = true
						consumer.conn.Close()
						consumer.es.staled <- consumer
						return
					}
				}
				idleTimer.Reset(es.idleTimeout)
			case <-idleTimer.C:
				consumer.conn.Close()
				consumer.es.staled <- consumer
				return
			}
		}
	}()

	return consumer, nil
}
