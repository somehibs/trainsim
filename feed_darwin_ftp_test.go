package trainsim

import (
	"fmt"
	"testing"
)

func TestDarwinFtp(t *testing.T) {
	df := NewDarwinFtp()
	ftpXml := df.Fetch()
	fmt.Printf("xml data lengths: time: %+v ref: %+v\n", len(ftpXml.Timetable.Journeys), len(ftpXml.Reference.Locations))
}
