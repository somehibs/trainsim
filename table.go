package trainsim

import (
	"fmt"
	"time"

	"git.circuitco.de/self/trainsim/xml"

	"github.com/jinzhu/gorm"
)

type Location struct {
	Id         int64 `gorm:"AUTO_INCREMENT"`
	OperatorId int64
	Ref        string `gorm:"primary_key"`
	Name       string
	AlphaCode  string
}

type Operator struct {
	Id   int64  `gorm:"AUTO_INCREMENT"`
	Ref  string `gorm:"primary_key"`
	Name string
	Url  string
}

type Reason struct {
	Namespace int32 `gorm:"primary_key"`
	Code      int64 `gorm:"primary_key"`
	Text      string
}

func ensureTypes(db *gorm.DB) {
	db.AutoMigrate(Journey{}, Operator{}, Reason{}, Location{}, CallingPoint{})
}

func (r Reason) FromXml(ns int32, o darwin.Reason) Reason {
	r.Namespace = ns
	r.Code = o.Code
	r.Text = o.Text
	return r
}

type Journey struct {
	Id               int64
	OperatorId       int64
	Origin           string // prefers OR unless OPOR is only field
	Destination      string // prefers DT unless OPDT is only field
	StartTime        time.Time
	EndTime          time.Time
	Cancelled        uint32 // 0 if not cancelled, otherwise code
	Uid              string
	TrainId          string
	StartDate        time.Time
	TrainCategory    string
	PassengerService bool // usually means only OPOR fields are present
	ServiceType      string
	QueueTrain       bool
	Charter          bool
}

func consumeJourney(db *gorm.DB, journey darwin.Journey) {
	if journey.Deleted {
		return
	}
	passengers := true
	if journey.IsPassengerSvc == "false" {
		passengers = false
	}
	origin := ""
	startTime := time.Unix(0, 0)
	if journey.Origin.Ref != "" {
		origin = journey.Origin.Ref
		startTime = CallingPoint{}.FromXml(journey.Origin, Journey{}, nil).GetOldest()
	} else {
		origin = journey.InternalOrigin.Ref
		startTime = CallingPoint{}.FromXml(journey.InternalOrigin, Journey{}, nil).GetOldest()
	}
	destination := ""
	endTime := time.Unix(0, 0)
	if journey.Destination.Ref != "" {
		destination = journey.Destination.Ref
		endTime = CallingPoint{}.FromXml(journey.Destination, Journey{}, nil).GetOldest()
	} else {
		destination = journey.InternalDestination.Ref
		endTime = CallingPoint{}.FromXml(journey.InternalDestination, Journey{}, nil).GetOldest()
	}

	cancelReason := uint32(0)
	if journey.Cancelled {
		for _, c := range journey.CancelReason {
			if c.LocationRef == "" {
				cancelReason = c.Code
			}
		}
	}
	started, _ := time.Parse("2006-01-02", journey.StartDate)
	j := Journey{
		Id:               journey.Id,
		OperatorId:       operatorIds[journey.OperatorRef],
		Origin:           origin,
		Destination:      destination,
		StartTime:        startTime,
		EndTime:          endTime,
		Cancelled:        cancelReason,
		Uid:              journey.Uid,
		TrainId:          journey.TrainId,
		StartDate:        started,
		TrainCategory:    journey.ServiceType,
		QueueTrain:       journey.QueueTrain,
		PassengerService: passengers,
		ServiceType:      journey.ServiceType,
		Charter:          journey.IsCharter,
	}
	inserted := db.Create(&j)
	dbj := *inserted.Value.(*Journey)
	// Now insert all calling points
	insertPoints(db, journey, dbj)
}

type CallingPoint struct {
	Id               int64 `gorm:"AUTO_INCREMENT"`
	Type             uint32
	Journey          int64
	Location         string `gorm:"index"` // reference to the reference xml with location data
	FalseLocation    string // some reason there's a second location on xml entities
	Cancelled        int32  // cancellation code or 0
	Activity         string // if the train is doing some other things at this calling point
	PlannedActivity  string // if the train... might be doing some other things? idk
	Platform         string // planned platform
	PublicArrival    time.Time
	PublicDeparture  time.Time
	WorkingArrival   time.Time
	WorkingPassed    time.Time
	WorkingDeparture time.Time
}

func (cp CallingPoint) ValidTime(t time.Time) bool {
	if t.Unix() == 0 {
		return false
	}
	return true
}

func (cp CallingPoint) GetOldest() time.Time {
	if cp.ValidTime(cp.PublicDeparture) {
		return cp.PublicDeparture
	} else if cp.ValidTime(cp.PublicArrival) {
		return cp.PublicArrival
	} else if cp.ValidTime(cp.WorkingDeparture) {
		return cp.WorkingDeparture
	} else if cp.ValidTime(cp.WorkingArrival) {
		return cp.WorkingArrival
	}
	return time.Now()
}

func (cp CallingPoint) After(other CallingPoint) bool {
	oldest := cp.GetOldest()
	otherOldest := other.GetOldest()
	return oldest.After(otherOldest)
}

var timeFormat = "15:04:05"

func MustParse(format, timeStr string) time.Time {
	if timeStr == "" {
		return time.Unix(0, 0)
	}
	if len(timeStr) == 5 {
		timeStr = timeStr + ":00"
	}
	time, err := time.Parse(format, timeStr)
	if err != nil {
		panic("Didn't understand time " + err.Error())
	}
	return time
}

func (cp CallingPoint) FromXml(p darwin.CallingPoint, journey Journey, cancellations map[string]darwin.JourneyCancelReason) CallingPoint {
	cp.Location = p.Ref
	cp.Journey = journey.Id
	cp.Activity = p.Activity
	cp.PlannedActivity = p.PlannedActivity
	cp.Platform = p.Platform
	cp.PublicArrival = MustParse(timeFormat, p.PublicArrival)
	cp.PublicDeparture = MustParse(timeFormat, p.PublicDeparture)
	cp.WorkingArrival = MustParse(timeFormat, p.WorkingArrival)
	cp.WorkingPassed = MustParse(timeFormat, p.WorkingPassed)
	cp.WorkingDeparture = MustParse(timeFormat, p.WorkingDeparture)
	return cp
}

func insertPoint(tx *gorm.DB, xml darwin.CallingPoint, cpType darwin.CallingPointType, dbJourney Journey) {
	point := CallingPoint{}.FromXml(xml, dbJourney, nil)
	point.Type = uint32(cpType)
	tx.Create(&point)
}

func insertPoints(db *gorm.DB, journey darwin.Journey, dbJourney Journey) {
	tx := db.Begin()
	if journey.Origin.Ref != "" {
		insertPoint(tx, journey.Origin, darwin.CallingPoint_Origin, dbJourney)
	} else if journey.InternalOrigin.Ref != "" {
		insertPoint(tx, journey.InternalOrigin, darwin.CallingPoint_OperationalOrigin, dbJourney)
	}
	for _, pass := range journey.PassingPoints {
		insertPoint(tx, pass, darwin.CallingPoint_Passing, dbJourney)
	}
	for _, calling := range journey.CallingPoints {
		insertPoint(tx, calling, darwin.CallingPoint_Intermediate, dbJourney)
	}
	tx.Commit()
}

func ConsumeXmlJourney(db *gorm.DB, data *NightlyXmlData) {
	for _, table := range []string{"journeys", "calling_points"} {
		db.Exec("DROP TABLE " + table)
	}
	fmt.Println("journey count", len(data.Timetable.Journeys))
	ensureTypes(db)
	CacheOperIds(db)
	// normalise each journey
	parsed := 0
	jl := len(data.Timetable.Journeys)
	for _, journey := range data.Timetable.Journeys {
		parsed += 1
		if parsed%1000 == 0 {
			pct := float64(parsed) / float64(jl)
			fmt.Println(parsed, "/", jl, "(", 100*pct, "%)")
		}
		consumeJourney(db, journey)
	}
}

var operatorIds = map[string]int64{}

func CacheOperIds(db *gorm.DB) {
	opers := make([]Operator, 0)
	db.Find(&opers)
	for _, oper := range opers {
		operatorIds[oper.Ref] = oper.Id
	}
}

func ConsumeXmlReference(db *gorm.DB, data *NightlyXmlData) {
	for _, table := range []string{"operators", "locations", "reasons"} {
		db.Exec("DROP TABLE " + table)
	}
	// Ensure types
	ensureTypes(db)
	// Process references
	// Start with cancellation and late reasons
	for _, reason := range data.Reference.LateReasons.Reasons {
		r := Reason{}.FromXml(1, reason)
		db.Create(&r)
	}

	for _, reason := range data.Reference.CancellationReasons.Reasons {
		r := Reason{}.FromXml(2, reason)
		db.Create(&r)
	}

	// Next, process Operators
	operatorIds = map[string]int64{}
	for _, oper := range data.Reference.Operators {
		op := Operator{Ref: oper.Ref, Name: oper.Name, Url: oper.URL}

		var existing Operator
		db.Where("ref = ?", oper.Ref).First(&existing)
		if existing.Id != 0 {
			operatorIds[existing.Ref] = existing.Id
			continue
		}
		inserted := db.Create(&op)
		existing = *inserted.Value.(*Operator)
		operatorIds[existing.Ref] = existing.Id
	}

	// Process locations
	for _, location := range data.Reference.Locations {
		loc := Location{
			OperatorId: operatorIds[location.OperatorRef],
			Ref:        location.Ref,
			AlphaCode:  location.AlphaCode,
			Name:       location.Name,
		}
		db.Create(&loc)
	}
}
