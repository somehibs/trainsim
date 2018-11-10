package trainsim

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"

	"git.circuitco.de/self/trainsim/xml"
)

type DarwinFtp struct {
	callback FetchCallback
}

var FtpDir = "xml/"

func NewDarwinFtp() DarwinFtp {
	df := DarwinFtp{}
	os.Mkdir(FtpDir, 0777)
	return df
}

func (df *DarwinFtp) Fetch() *NightlyXmlData {
	// Check if the files on disk are up to date or not
	if !df.ftpFetched() {
		fmt.Println("fetching from ftp")
		df.fetchFtp()
	}
	if !df.xmlExtracted() {
		fmt.Println("gunzipping xml")
		df.extractXml()
	}
	// Always parse from XML, unless cached
	ret := df.parseXml()
	fmt.Println("fetch complete")
	return ret
}

func (df *DarwinFtp) xmlExtracted() bool {
	return df.fileCheck("xml")
}

func (df *DarwinFtp) ftpFetched() bool {
	return df.fileCheck("gz")
}

func (df *DarwinFtp) fileCheck(ext string) bool {
	// Check the filesystem
	f, e := ioutil.ReadDir(FtpDir)
	if e != nil {
		panic(fmt.Sprintf("Could not read dir %s", FtpDir))
	}
	requiredFiles := []string{
		"ref_v3",
		"v8",
	}
	found := 0
	for _, file := range f {
		n := file.Name()
		if !strings.HasSuffix(n, ext) {
			continue
		}
		matched := false
		for _, req := range requiredFiles {
			if strings.Contains(n, req) {
				matched = true
				break
			}
		}
		if !matched {
			fmt.Printf("Ignoring file %s\n", n)
			continue
		}
		xmln := df.parseFilename(n)
		if xmln.today {
			found += 1
		}
	}
	return len(requiredFiles) == found
}

type XmlName struct {
	time  time.Time
	name  string
	today bool
}

var fileTime = "20060102150405"

func (df *DarwinFtp) parseFilename(name string) XmlName {
	xn := XmlName{}
	prefix := strings.SplitN(name, "_", 2)
	if len(prefix) > 1 {
		ft, err := time.Parse(fileTime, prefix[0])
		if err != nil {
			panic(fmt.Sprintf("Unexpected unmatched datekey %s", ft))
		}
		nt := time.Now()
		if ft.Day() == nt.Day() && ft.Month() == nt.Month() && nt.Year() == ft.Year() {
			xn.today = true
		}
		xn.time = ft
		fname := strings.Split(prefix[1], ".")
		xn.name = fname[0]
	}
	return xn
}

func (df *DarwinFtp) extractXml() {
	filepath.Walk(FtpDir, func(root string, info os.FileInfo, err error) error {
		if df.fileOkay(info, ".gz") != nil {
			GunzipFile(FtpDir+info.Name(), FtpDir+info.Name()[:len(info.Name())-3])
		}
		return nil
	})
}

func (df *DarwinFtp) fileOkay(info os.FileInfo, ext string) *XmlName {
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(info.Name(), ext) {
		return nil
	}
	xn := df.parseFilename(info.Name())
	if xn.name == "" {
		return nil
	}
	return &xn
}

type NightlyXmlData struct {
	Timetable darwin.PportTimetable
	Reference darwin.PportTimetableRef
}

func (df *DarwinFtp) parseXml() *NightlyXmlData {
	var x NightlyXmlData
	filepath.Walk(FtpDir, func(root string, info os.FileInfo, err error) error {
		xn := df.fileOkay(info, ".xml")
		if xn != nil {
			fmt.Println(xn.name)
			if !xn.today {
				fmt.Printf("Refusing to parse stale xml %s", xn)
				return nil
			}
			f, e := os.Open(FtpDir + info.Name())
			if e != nil {
				return e
			}
			b, e := ioutil.ReadAll(f)
			if e != nil {
				return e
			}
			if xn.name == "ref_v3" {
				fmt.Println("Unmarshalling reference")
				xml.Unmarshal(b, &x.Reference)
			} else if xn.name == "v8" {
				fmt.Println("Unmarshalling timetable")
				xml.Unmarshal(b, &x.Timetable)
			}
		}
		return nil
	})
	if len(x.Timetable.Journeys) == 0 || len(x.Timetable.Associations) == 0 || len(x.Reference.Locations) == 0 {
		panic("loaded timetable or reference appears invalid")
	}
	return &x
}

func (df *DarwinFtp) fetchFtp() {
	// Kick off an FTP connection to retrieve the file to somewhere local for further parsing by the FetchItem
	config := GetDarwinFtp()
	conn, e := ftp.Connect(config.Hostname)
	if e != nil {
		panic("Couldn't connect: " + e.Error())
	}
	defer conn.Quit()
	e = conn.Login(config.Username, config.Password)
	if e != nil {
		panic("Couldn't login: " + e.Error())
	}
	entries, e := conn.List("")
	if e != nil {
		panic("Couldn't list FTP: " + e.Error())
	}
	df.FetchXml(conn, entries)
}

func (df *DarwinFtp) FetchXml(conn *ftp.ServerConn, entries []*ftp.Entry) {
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name, ".xml.gz") {
			// Fetch all the XML files and store them on disk under the same name
			_, e := os.Open(FtpDir + entry.Name)
			if e == nil {
				fmt.Println("this xml file already exists")
				continue
			}
			f, e := os.Create(FtpDir + entry.Name)
			if e != nil {
				panic("Couldn't open: " + FtpDir + entry.Name)
			}
			r, e := conn.Retr(entry.Name)
			if e != nil {
				fmt.Println("could not fetch: " + e.Error())
				r, e = conn.Retr(entry.Name)
				r.Close()
				if e != nil {
					panic(fmt.Sprintf("Could not retrieve %s (%s)", entry.Name, e.Error()))
				}
			}
			io.Copy(f, r)
			r.Close()
			fmt.Printf("Entry: %+v\n", entry)
		}
	}
}
