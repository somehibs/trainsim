package darwin

import "encoding/xml"

// PportTimetable was generated 2018-11-10 01:15:22
type PportTimetable struct {
	TimetableID string        `xml:"timetableID,attr"`
	Journey     []Journey     `xml:"Journey"`
	Association []Association `xml:"Association"`

	XMLName xml.Name `xml:"PportTimetable"`
	Xsd     string   `xml:"xsd,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
}

// manually resolved xsd and xml structs
type PassingPoint struct {
	Tiploc     string `xml:"tpl,attr"`
	PassedAt   string `xml:"wtp,attr"`
	Platform   string `xml:"plat,attr"`
	Action     string `xml:"act,attr"`
	Cancelled  string `xml:"can,attr"`
	RouteDelay string `xml:"rdelay,attr"`
}

type CallingPoint struct {
	Tiploc              string `xml:"tpl,attr"`
	Activity            string `xml:"act,attr"`
	PlannedActivity     string `xml:"planAct,attr"`
	Platform            string `xml:"plat,attr"`
	PublicArrival       string `xml:"pta,attr"`
	PublicDeparture     string `xml:"ptd,attr"`
	WorkingArrival      string `xml:"wta,attr"`
	WorkingDeparture    string `xml:"wtd,attr"`
	Cancelled           string `xml:"can,attr"`
	FalseDestinationTip string `xml:"fd,attr"`
	RouteDelay          string `xml:"rdelay,attr"`
}

type Journey struct {
	Rid                      string         `xml:"rid,attr"`
	Uid                      string         `xml:"uid,attr"`
	TrainId                  string         `xml:"trainId,attr"`
	Ssd                      string         `xml:"ssd,attr"`
	Toc                      string         `xml:"toc,attr"`
	TrainCat                 string         `xml:"trainCat,attr"`
	IsPassengerSvc           string         `xml:"isPassengerSvc,attr"`
	Status                   string         `xml:"status,attr"`
	Qtrain                   string         `xml:"qtrain,attr"`
	Deleted                  string         `xml:"deleted,attr"`
	Can                      string         `xml:"can,attr"`
	IsCharter                string         `xml:"isCharter,attr"`
	OriginPoint              []CallingPoint `xml:"OR"`
	PassingPoints            []PassingPoint `xml:"PP"`
	CallingPoints            []CallingPoint `xml:"IP"`
	DestinationPoint         []CallingPoint `xml:"DT"`
	InternalOriginPoint      CallingPoint   `xml:"OPOR"`
	InternalCallingPoint     []CallingPoint `xml:"OPIP"`
	InternalDestinationPoint CallingPoint   `xml:"OPDT"`
	CancelReason             struct {
		Code   int    `xml:",chardata"`
		Tiploc string `xml:"tiploc,attr"`
	} `xml:"cancelReason"`
}

type Association struct {
	Text     string `xml:",chardata"`
	Tiploc   string `xml:"tiploc,attr"`
	Category string `xml:"category,attr"`
	Main     struct {
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
	Text    string `xml:",chardata"`
	Tpl     string `xml:"tpl,attr"`
	Locname string `xml:"locname,attr"`
	Crs     string `xml:"crs,attr"`
	Toc     string `xml:"toc,attr"`
}

type TermsLink struct {
	Text    string `xml:",chardata"`
	Toc     string `xml:"toc,attr"`
	Tocname string `xml:"tocname,attr"`
	URL     string `xml:"url,attr"`
}

type Reason struct {
	Code int    `xml:"code,attr"`
	Text string `xml:"reasontext,attr"`
}

type PportTimetableRef struct {
	TimetableId string      `xml:"timetableId,attr"`
	XMLName     xml.Name    `xml:"PportTimetableRef"`
	Text        string      `xml:",chardata"`
	Xsd         string      `xml:"xsd,attr"`
	Xsi         string      `xml:"xsi,attr"`
	Xmlns       string      `xml:"xmlns,attr"`
	Locations   []Location  `xml:"LocationRef"`
	TermsLinks  []TermsLink `xml:"TocRef"`
	Reasons     struct {
		ReasonReasons []Reason `xml:"Reason"`
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
