package snowflake

import (
	"time"
)

const (
	COUNTER_BIT_SIZE int = 14
	THREAD_CAP       int = 1 << COUNTER_BIT_SIZE
	THREAD_BIT_SIZE  int = 4
)

type Generator struct {
	Epoch        int64
	GeneratedIDs []chan int64
	ThreadCount  int
	done         chan bool
}

func NewWithDefaultEpoch(threadCount int) *Generator {
	generatedIDs := make([]chan int64, threadCount)
	for i := 0; i < threadCount; i++ {
		generatedIDs[i] = make(chan int64, 40000000)
	}
	return &Generator{
		Epoch:        time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC).UnixMilli(),
		GeneratedIDs: generatedIDs,
		ThreadCount:  len(generatedIDs),
		done:         make(chan bool),
	}
}

func (g Generator) generateID(threadId int) {
	prevTime := time.Now().UnixMilli()
	counter := 0
	for {
		select {
		case <-g.done:
			close(g.GeneratedIDs[threadId])
			return
		default:
			currentTime := time.Now().UnixMilli()
			if currentTime != prevTime {
				// Not the same millisecond.
				prevTime = currentTime
				counter = 0
			} else if counter >= THREAD_CAP {
				// Prevent overflow the 14 bits counter in the same millisecond.
				time.Sleep(10 * time.Microsecond)
			}
			timeSinceEpoch := currentTime - g.Epoch
			// timestamp,    thread,       counter
			//            [  4 bits  ]  [  14 bits  ]
			id := timeSinceEpoch << 18
			id |= int64(threadId) << COUNTER_BIT_SIZE
			id |= int64(counter)
			g.GeneratedIDs[threadId] <- id
			counter++
		}
	}
}

func (g Generator) Start() {
	for i := 0; i < g.ThreadCount; i++ {
		go g.generateID(i)
	}
}

func (g Generator) Done() {
	g.done <- true
}
