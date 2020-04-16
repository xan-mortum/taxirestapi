package components

import (
	"math/rand"
	"time"
)

const countOfRequests = 50
const aCharIndex = 97
const zCharIndex = 122
const timeToChangeRequest = 200

type DB struct {
	requests      []string
	requestsCalls *Counter

	done           chan bool
	getRequestIn   chan bool
	getRequestOut  chan string
	getRequestsIn  chan bool
	getRequestsOut chan map[string]int
}

func NewDB() *DB {
	return &DB{
		requests:       make([]string, countOfRequests),
		requestsCalls:  NewCounter(),
		done:           make(chan bool),
		getRequestIn:   make(chan bool),
		getRequestOut:  make(chan string),
		getRequestsIn:  make(chan bool),
		getRequestsOut: make(chan map[string]int),
	}
}

func (db *DB) Generate() {
	i := 0
	for i < countOfRequests {
		request := db.getRandomRequestName()
		db.requests[i] = request
		db.requestsCalls.Store(request, 0)
		i++
	}
}

func (db *DB) Start() {
	go func() {
		for {
			select {
			case <-db.getRequestsIn:
				db.getRequestsOut <- db.requestsCalls.List()
			case <-db.getRequestIn:
				db.getRequestOut <- db.requests[rand.Intn(countOfRequests-1)]
			case <-db.done:
				return
			case <-time.After(timeToChangeRequest * time.Millisecond):
				request := db.getRandomRequestName()
				db.requests[rand.Intn(countOfRequests-1)] = request
				db.requestsCalls.Inc(request)
			}
		}
	}()
}

func (db *DB) Stop() {
	go func() {
		db.done <- true
	}()
}

func (db *DB) GetRequest() string {
	go func() {
		db.getRequestIn <- true
	}()
	request := <-db.getRequestOut
	db.requestsCalls.Inc(request)
	return request
}

func (db *DB) GetRequests() map[string]int {
	go func() {
		db.getRequestsIn <- true
	}()
	requests := <-db.getRequestsOut
	result := make(map[string]int)
	for request, calls := range requests {
		if calls < 1 {
			continue
		}
		result[request] = calls
	}
	return result
}

func (db *DB) getRandomRequestName() string {
	min := aCharIndex
	max := zCharIndex
	return string(rune(rand.Intn(max-min)+min)) + string(rune(rand.Intn(max-min)+min))
}
