package store

import (
	"log"
	"time"
)

type Store struct {
	store  map[string]string
	logger *log.Logger
}

func New(logger *log.Logger) *Store {
	return &Store{
		store:  make(map[string]string),
		logger: logger,
	}
}

func (s *Store) Add(key, value string) {
	s.logger.Println("adding key/value to store")
	time.Sleep(time.Second * 2)
	s.store[key] = value
	s.logger.Println("added key/value to store")
}

func (s *Store) Get(key string) string {
	return s.store[key]
}

func (s *Store) LogAll() {
	for k, v := range s.store {
		s.logger.Printf("%s: %s\n", k, v)
	}
}
