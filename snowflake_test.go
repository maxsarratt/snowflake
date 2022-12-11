package snowflake_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/maxsarratt/snowflake"
)

func TestBenchmark(t *testing.T) {
	threadCount := 1

	startTime := time.Now().UnixMicro()
	g := snowflake.NewWithDefaultEpoch(threadCount)

	g.Start()
	time.Sleep(100 * time.Millisecond)
	g.Done()

	endTime := time.Now().UnixMicro()

	count := 0
	for i := 0; i < threadCount; i++ {
		count += len(g.GeneratedIDs[i])
	}
	duration := (endTime - startTime) / 1000.0

	fmt.Printf("\n %f \n", float64(count)/float64(duration))
}
