package main

import (
	"flag"

	"git.circuitco.de/self/trainsim"
)

func parseXml() {
	ftp := trainsim.NewDarwinFtp()
	xml := ftp.Fetch()
	db := trainsim.NewDb()
	trainsim.ConsumeXmlJourney(db, xml)
	db.Close()
}

func main() {
	shouldParse := flag.Bool("parse", false, "parse from xml and populate database")
	flag.Parse()
	if *shouldParse {
		parseXml()
	}
	s := trainsim.NewSelect()
	defer s.Close()
	s.PlainRoute(s.LocationByCode("LDS"), s.LocationByCode("MAN"))
	s.PlainRoute(s.LocationByCode("MAN"), s.LocationByCode("LDS"))
}
