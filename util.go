package trainsim

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func PerfStart() time.Time {
	return time.Now()
}

func PerfEnd(component string, start time.Time) {
	fmt.Printf("perf: %s: %f\n", component, time.Now().Sub(start).Seconds()*1000) // search took n ms
}

func GunzipFile(file, extracted string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()
	gzf, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	ex, err := os.Create(extracted)
	if err != nil {
		return false
	}
	defer ex.Close()
	io.Copy(ex, gzf)
	return true
}

func ReadJson(file string, item interface{}) error {
	f, e := os.Open(ConfDir + file)
	if e != nil {
		return e
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, item)
}
