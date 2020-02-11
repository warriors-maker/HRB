package Server

import (
	"sync"
	"time"
)


type StatsCounter struct {
	broadCastCounter int
	m sync.Mutex
}

func (counter StatsCounter) increment(){
	counter.m.Lock()
	counter.broadCastCounter += 1
	counter.m.Unlock()
}

func (counter StatsCounter) getCount() int{
	var count int
	counter.m.Lock()
	count = counter.broadCastCounter
	counter.m.Unlock()
	return count
}

func startThroughput() {
	time.Sleep(1* time.Minute)
	//Write to File
}