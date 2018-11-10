package trainsim

// API for retrieving data from NRDP
// Calls out to other APIs for retrieving and translating HSP and DFTP
// Subscribes to feeds published to RTF
var ConfDir = "conf/"

type FetchCallback func(FeedItem)

type FeedItem interface {
	Raw() []byte
	Parse() interface{}
}

type Feed interface {
	Fetch(FetchCallback)
}

// HSP config
type HspConfig struct {
	Hostname          string
	ServiceDetailUrl  string
	ServiceMetricsUrl string
}

func GetHspConfig() HspConfig {
	return HspConfig{
		"hsp-prod.rockshore.net",
		"/api/v1/serviceDetails",
		"/api/v1/serviceMetrics",
	}
}

type DarwinStompConfig struct {
	Queue    string
	Username string
	Password string
}

func GetDarwinStomp() DarwinStompConfig {
	d := DarwinStompConfig{
		"",
		"d3user",
		"d3password",
	}
	e := ReadJson("darwin-stomp.json", &d)
	if e != nil {
		panic("Failed to read stomp config: " + e.Error())
	}
	return d
}

// Darwin FTP config
type DarwinFtpConfig struct {
	Hostname string
	Username string
	Password string
}

func GetDarwinFtp() DarwinFtpConfig {
	conf := DarwinFtpConfig{
		"datafeeds.nationalrail.co.uk:21",
		"ftpuser",
		"",
	}
	e := ReadJson("darwin-ftp.json", &conf)
	if e != nil {
		panic("failed to read config from darwin-ftp.json")
	}
	return conf
}
