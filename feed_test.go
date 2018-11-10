package trainsim

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	//	f := NewFeeds()
	fmt.Println(GetHspConfig())
	fmt.Println(GetDarwinFtp())
	fmt.Println(GetDarwinStomp())
}
