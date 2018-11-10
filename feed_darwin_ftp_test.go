package trainsim

import (
	"testing"
)

func TestDarwinFtp(t *testing.T) {
	df := NewDarwinFtp()
	df.Fetch()
}
