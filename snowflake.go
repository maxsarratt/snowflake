package snowflake

import (
	"sync"
	"time"
)

const (
	COUNTER_BITS int8 = 14
	THREAD_CAP   int  = 1 << COUNTER_BITS
)

type Generator struct {
	sync.Mutex
	// Time in milliseconds.
	epoch    int64
	prevTime int64
	counter  int
}

func NewWithDefaultEpoch() *Generator {
	return &Generator{
		epoch:    time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC).UnixMilli(),
		prevTime: time.Now().UnixMilli(),
		counter:  0,
	}
}

func (g *Generator) Generate() int64 {
	g.Lock()
	defer g.Unlock()

	currentTime := time.Now().UnixMilli()

	if currentTime == g.prevTime {
		if g.counter >= THREAD_CAP {
			// Prevent overflow the 14 bits counter in the same millisecond.
			for currentTime == g.prevTime {
				currentTime = time.Now().UnixMilli()
			}
			g.counter = 0
		}
	} else {
		g.counter = 0
	}

	timeSinceEpoch := currentTime - g.epoch
	// timestamp,    counter
	//            [  14 bits  ]
	id := timeSinceEpoch << COUNTER_BITS
	id |= int64(g.counter)

	g.prevTime = currentTime
	g.counter += 1

	return id
}
