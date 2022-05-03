package server

import (
	"log"
	"os"
	"testing"

	"github.com/tomknightdev/tcp/store"
)

func TestHandleRequestAdd(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		l := log.New(os.Stdout, "tcp-server-test", log.LstdFlags|log.Lshortfile)
		s := store.NewMemStore(l)
		tcp := NewTCP(l, s)

		request := []string{"add", "key", "value"}
		successResponse := tcp.handleRequest(request)

		if successResponse != "added" {
			t.Error("did not receive added response")
		}

		responseChan := make(chan store.Response)
		tcp.store.Get("key", responseChan)
		message := <-responseChan
		if message.Message != "value" {
			t.Errorf("got %s wanted %s", message.Message, "value")
		}
	})
}
