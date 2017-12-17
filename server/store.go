package main

import (
	"bytes"
	"sync"
	"time"

	"github.com/kataras/golog"
)

type commandStore struct {
	m    sync.Mutex
	dict map[string]commandStoreItem

	readTimeout     time.Duration
	writeTimeout    time.Duration
	creationTimeout time.Duration
}

type commandStoreItem struct {
	command       *command
	creationDate  time.Time
	lastWriteDate time.Time
	lastReadDate  time.Time
	stdin         *bytes.Buffer
	stdout        *bytes.Buffer
	stderr        *bytes.Buffer
}

func newCommandStore() *commandStore {
	return &commandStore{
		dict:            make(map[string]commandStoreItem),
		readTimeout:     time.Second * 60,
		writeTimeout:    time.Second * 120,
		creationTimeout: 0,
	}
}

func (s *commandStore) Get(key string) commandStoreItem {
	s.m.Lock()
	defer s.m.Unlock()

	value, ok := s.dict[key]

	// to force copy of the item for update lastReadDate
	func(value commandStoreItem) {
		if ok {
			value.lastReadDate = time.Now()
			s.dict[key] = value
		}
	}(value)

	return value
}

func (s *commandStore) Set(key string, value commandStoreItem) commandStoreItem {
	s.m.Lock()
	defer s.m.Unlock()

	if _, found := s.dict[key]; !found && value.creationDate.IsZero() {
		value.creationDate = time.Now()
	}
	value.lastWriteDate = time.Now()
	s.dict[key] = value
	return value
}

func (s *commandStore) Size() int {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.dict)
}

func (s *commandStore) clear(log golog.Logger) {
	s.m.Lock()
	defer s.m.Unlock()

	now := time.Now()

	for k, v := range s.dict {
		if s.readTimeout != 0 && v.lastReadDate.Add(s.readTimeout).Before(now) {
			log.Debug("commandStore.clear: id=%+v evicted due to read timeout", k)
			delete(s.dict, k)
			continue
		}
		if s.writeTimeout != 0 && v.lastWriteDate.Add(s.writeTimeout).Before(now) {
			log.Debug("commandStore.clear: id=%+v evicted due to write timeout", k)
			delete(s.dict, k)
			continue
		}
		if s.creationTimeout != 0 && v.creationDate.Add(s.creationTimeout).Before(now) {
			log.Debug("commandStore.clear: id=%+v evicted due to creation timeout", k)
			delete(s.dict, k)
			continue
		}
	}
}
