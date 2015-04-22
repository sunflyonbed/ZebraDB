package main

type xmlConfig struct {
	Log xmlLog `xml:"log"`
	DB  xmlDB  `xml:"db"`
}

type xmlLog struct {
	Config string `xml:"config"`
}

type xmlDB struct {
	Path         string `xml:"path"`
	MonitorAddr  string `xml:"monitor_addr"`
	MonitorIndex int    `xml:"monitor_index"`
	ListenAddr   string `xml:"listen_addr"`
}
