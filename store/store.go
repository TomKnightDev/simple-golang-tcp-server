package store

import (
	"log"
)

type Store interface {
	Add(string, string, chan Response)
	Get(string, chan Response)
	LogAll()
}

type MemStore struct {
	store  map[string]string
	logger *log.Logger
}

var requestChan = make(chan request)

type request struct {
	reqType      string
	key          string
	value        string
	responseChan chan Response
}

type Response struct {
	Message string
	Err     error
}

func NewMemStore(logger *log.Logger) *MemStore {
	s := &MemStore{
		store:  make(map[string]string),
		logger: logger,
	}

	go func() {
		for {
			request := <-requestChan

			if request.reqType == "get" {
				request.responseChan <- Response{s.store[request.key], nil}
			} else if request.reqType == "add" {
				s.store[request.key] = request.value
				request.responseChan <- Response{"added", nil}
			}
		}
	}()

	return s
}

func (s *MemStore) Add(key string, value string, responseChan chan Response) {
	go func() {
		requestChan <- request{"add", key, "value", responseChan}
	}()
}

func (s *MemStore) Get(key string, responseChan chan Response) {
	go func() {
		requestChan <- request{"get", key, "", responseChan}
	}()
}

func (s *MemStore) Delete(key string) {
	delete(s.store, key)
}

func (s *MemStore) LogAll() {
	for k, v := range s.store {
		s.logger.Printf("%s: %s\n", k, v)
	}
}
