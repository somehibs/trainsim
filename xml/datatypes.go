package darwin

import "encoding/xml"

var Activities = map[string]string{
	"A":  "Stops or shunts for other trains to pass",
	"AE": "(De)couple assisting locomotive",
	"AX": "Shows as X on arrival", // I don't know what this means either
	"BL": "Stops for 'banking locomotive",
	"C":  "Stops to change trainmen",
	"D":  "Stops to set down passengers",
	"-D": "Stops to detach vehicles",
	"E":  "Stops for examination",
	"G":  "National Rail Timetable data to add",                  // no clue
	"H":  "Notional activity to prevent WTT timing column merge", // ??
	"HH": "As H where a third colum is involved",                 // ...k
	"K":  "Passenger count point",
	"KC": "Ticket collection and examination point",
	"KE": "Ticket examination point",
	"KF": "Ticket examination first class",
	"KS": "Selective ticket examination",
	"L":  "Locomotive change",
	"N":  "Stop not advertised", // unexpected stop
	"OP": "Operational reasons (undefined)",
	"OR": "Locomotive on rear",
	"PR": "Propelling between points",
	"R":  "Stops when required",
	"RM": "Reversing movement, or driver changes ends",
	"RR": "Stops for locomotive to run around train",
	"S":  "Stops only for rail personel",
	"T":  "Passenger (dis)embark",
	"-T": "Attach/detail",
	"TB": "Train begins",
	"TF": "Train finishes",
	"TS": "Detail Consist for TOPS Direct",
	"TW": "Stops or at pass for tablet, staff or token",
	"U":  "Stops to attach vehicles",
	"W":  "Stops for watering of the coaches",
	"X":  "Passes another train at crossing point on single line",
}

// PportTimetable was generated 2018-11-10 01:15:22
type PportTimetable struct {
	TimetableID  string        `xml:"timetableID,attr"`
	Journeys     []Journey     `xml:"Journey"`
	Associations []Association `xml:"Association"`

	XMLName xml.Name `xml:"PportTimetable"`
	Xsd     string   `xml:"xsd,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
}

const (
	// Undef
	CallingPoint_Undefined CallingPointType = 0

	// Of public interest
	CallingPoint_Origin       CallingPointType = 1
	CallingPoint_Passing      CallingPointType = 2
	CallingPoint_Destination  CallingPointType = 3
	CallingPoint_Intermediate CallingPointType = 4

	// Of internal interest
	CallingPoint_OperationalMin          CallingPointType = 10
	CallingPoint_OperationalOrigin       CallingPointType = 11
	CallingPoint_OperationalIntermediate CallingPointType = 12
	CallingPoint_OperationalDestination  CallingPointType = 13
)

type CallingPointType int

type CallingPoint struct {
	Ref                 string `xml:"tpl,attr"`
	Activity            string `xml:"act,attr"`
	PlannedActivity     string `xml:"planAct,attr"`
	Platform            string `xml:"plat,attr"`
	PublicArrival       string `xml:"pta,attr"`
	PublicDeparture     string `xml:"ptd,attr"`
	WorkingArrival      string `xml:"wta,attr"`
	WorkingDeparture    string `xml:"wtd,attr"`
	WorkingPassed       string `xml:"wtp,attr"` // for passing points only
	Cancelled           string `xml:"can,attr"`
	FalseDestinationTip string `xml:"fd,attr"`
	RouteDelay          string `xml:"rdelay,attr"`
}

type Journey struct {
	Id                   int64                 `xml:"rid,attr"`            // always present
	Uid                  string                `xml:"uid,attr"`            // always present
	TrainId              string                `xml:"trainId,attr"`        // always present
	StartDate            string                `xml:"ssd,attr"`            // always present
	OperatorRef          string                `xml:"toc,attr"`            // always present
	TrainCategory        string                `xml:"trainCat,attr"`       // sometimes present, default PP or something
	IsPassengerSvc       string                `xml:"isPassengerSvc,attr"` // ispassengerservice is true by default
	ServiceType          string                `xml:"status,attr"`         // sometimes 1. usually p
	QueueTrain           bool                  `xml:"qtrain,attr"`         // qtrain is 'runs as required' i.e. doesn't run unless there's a queue
	Deleted              bool                  `xml:"deleted,attr"`        // does actually exist
	Cancelled            bool                  `xml:"can,attr"`            // cancelled journey
	IsCharter            bool                  `xml:"isCharter,attr"`      // not seen it
	Origin               CallingPoint          `xml:"OR"`
	PassingPoints        []CallingPoint        `xml:"PP"` // passingpoint is a calling point with less data and a working time passed field
	CallingPoints        []CallingPoint        `xml:"IP"`
	Destination          CallingPoint          `xml:"DT"`
	InternalOrigin       CallingPoint          `xml:"OPOR"`
	InternalCallingPoint []CallingPoint        `xml:"OPIP"`
	InternalDestination  CallingPoint          `xml:"OPDT"`
	CancelReason         []JourneyCancelReason `xml:"cancelReason"`
}

type JourneyCancelReason struct {
	Code        uint32 `xml:",chardata"`
	LocationRef string `xml:"tiploc,attr"`
}

type Association struct {
	Text        string `xml:",chardata"`
	LocationRef string `xml:"tiploc,attr"`
	Category    string `xml:"category,attr"`
	Main        struct {
		Text string `xml:",chardata"`
		Rid  string `xml:"rid,attr"`
		Wta  string `xml:"wta,attr"`
		Pta  string `xml:"pta,attr"`
		Wtd  string `xml:"wtd,attr"`
		Ptd  string `xml:"ptd,attr"`
	} `xml:"main"`
	Assoc struct {
		Text string `xml:",chardata"`
		Rid  string `xml:"rid,attr"`
		Wtd  string `xml:"wtd,attr"`
		Ptd  string `xml:"ptd,attr"`
		Wta  string `xml:"wta,attr"`
		Pta  string `xml:"pta,attr"`
	} `xml:"assoc"`
}

// PportTimetableRef was generated 2018-11-10 01:15:22
type Location struct {
	Ref         string `xml:"tpl,attr"`
	Name        string `xml:"locname,attr"`
	AlphaCode   string `xml:"crs,attr"` // three character 'NOT' 'SHF' style code
	OperatorRef string `xml:"toc,attr"`
}

type Operator struct {
	Ref  string `xml:"toc,attr"`
	Name string `xml:"tocname,attr"`
	URL  string `xml:"url,attr"`
}

type Reason struct {
	Code int64  `xml:"code,attr"`
	Text string `xml:"reasontext,attr"`
}

type PportTimetableRef struct {
	TimetableId string     `xml:"timetableId,attr"`
	XMLName     xml.Name   `xml:"PportTimetableRef"`
	Text        string     `xml:",chardata"`
	Xsd         string     `xml:"xsd,attr"`
	Xsi         string     `xml:"xsi,attr"`
	Xmlns       string     `xml:"xmlns,attr"`
	Locations   []Location `xml:"LocationRef"`
	Operators   []Operator `xml:"TocRef"`
	LateReasons struct {
		Reasons []Reason `xml:"Reason"`
	} `xml:"LateRunningReasons"`
	CancellationReasons struct {
		Reasons []Reason `xml:"Reason"`
	} `xml:"CancellationReasons"`
	Via []struct {
		Target       string `xml:"at,attr"`
		Destination  string `xml:"dest,attr"`
		Text         string `xml:"viatext,attr"`
		ValidFrom    string `xml:"loc1,attr"`
		ValidFromTwo string `xml:"loc2,attr"`
	} `xml:"Via"`
	CisCodes []struct {
		Code string `xml:"code,attr"`
		Name string `xml:"name,attr"`
	} `xml:"CISSource"`
}
