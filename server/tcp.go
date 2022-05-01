package server

import (
	"bufio"
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tomknightdev/tcp/store"
)

type TCP struct {
	logger   *log.Logger
	store    store.Store
	listener net.Listener
	wg       sync.WaitGroup
}

func NewTCP(logger *log.Logger, store store.Store) *TCP {
	return &TCP{
		logger: logger,
		store:  store,
	}
}

func StartTCPServer(ctx context.Context, logger *log.Logger, store store.Store) {
	t := NewTCP(logger, store)

	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		t.logger.Println("error starting tcp server ", err)
	}
	t.listener = l

	t.wg.Add(1)
	go t.serve(ctx)

	<-ctx.Done()
	logger.Println("closing listener")
	t.listener.Close()
	logger.Println("waiting for group")
	t.wg.Wait()
	t.store.LogAll()
	logger.Println("server closed")
}

func (t *TCP) serve(ctx context.Context) {
	defer t.wg.Done()

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				t.logger.Println("error accepting tcp connection ", err)
			}
		} else {
			t.wg.Add(1)
			go func() {
				t.handleConnection(ctx, conn)
				t.wg.Done()
			}()
		}
	}
}

func (t *TCP) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetDeadline(time.Now().Add(time.Millisecond * 200))
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				} else {
					t.logger.Println("error reading message from client ", err)
					return
				}
			} else {
				conn.Write([]byte(data))
				t.store.Add(time.Now().String(), "value")
			}
		}
	}
}
