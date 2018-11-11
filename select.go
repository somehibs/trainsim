package trainsim

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Select struct {
	db *gorm.DB
}

func NewSelect() Select {
	sel := Select{}
	sel.db = NewDb()
	return sel
}

func (s Select) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s Select) LocationByRef(ref string) Location {
	var loc Location
	s.db.Where("ref = ?", ref).First(&loc)
	return loc
}

func (s Select) LocationByCode(crs string) Location {
	var loc Location
	s.db.Where("alpha_code = ?", crs).First(&loc)
	return loc
}

// Returns a list of journey ids
func (s Select) PlainRoute(start, end Location) []int64 {
	totalTime := PerfStart()
	//cps := make([]CallingPoint, 0)
	fmt.Println("Looking for route which contains", start.Name, "and", end.Name)
	startPoints := make([]CallingPoint, 0)
	endPoints := make([]CallingPoint, 0)
	a := PerfStart()
	s.db.Where("location = ?", start.Ref).Find(&startPoints)
	PerfEnd("select start", a)
	b := PerfStart()
	s.db.Where("location = ?", end.Ref).Find(&endPoints)
	PerfEnd("select ends", b)
	journeyPointCount := map[int64]CallingPoint{}
	c := PerfStart()
	total := 0
	for _, point := range startPoints {
		total += 1
		journeyPointCount[point.Journey] = point
	}
	for _, point := range endPoints {
		total += 1
		startPoint := journeyPointCount[point.Journey]
		if startPoint.Location == "" {
			continue
		} else if startPoint.After(point) {
			journeyPointCount[point.Journey] = CallingPoint{}
			continue
		}
	}
	j := []int64{}
	for k, v := range journeyPointCount {
		if v.Location == "" { // nonexisting point
			continue
		}
		j = append(j, k)
	}
	fmt.Printf("journeys with both %s and %s: %d (proc %d)\n", start.Name, end.Name, len(j), total)
	fmt.Printf("")
	PerfEnd("parse", c)
	PerfEnd("total query", totalTime)
	return j
}
