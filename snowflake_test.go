package snowflake_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/maxsarratt/snowflake"
)

func TestBenchmark(t *testing.T) {
	g := snowflake.NewWithDefaultEpoch()

	startTime := time.Now().UnixMicro()
	count := 0
	timeout := time.After(100 * time.Millisecond)

loop:
	for {
		select {
		case <-timeout:
			fmt.Println("There's no more time to this. Exiting!")
			break loop
		default:
			g.Generate()
			count++
		}
	}

	endTime := time.Now().UnixMicro()

	duration := (endTime - startTime) / 1000.0

	fmt.Printf("\n %f \n", float64(count)/float64(duration))
}

func BenchmarkGenerate(b *testing.B) {
	g := snowflake.NewWithDefaultEpoch()
	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
