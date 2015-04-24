package main

type xmlConfig struct {
	Log     xmlLog     `xml:"log"`
	Server  xmlServer  `xml:"server"`
	LevelDB xmlLevelDB `xml:"leveldb"`
}

type xmlLog struct {
	Config string `xml:"config"`
}

type xmlServer struct {
	Path         string `xml:"path"`
	MonitorAddr  string `xml:"monitor_addr"`
	MonitorIndex int    `xml:"monitor_index"`
	ListenAddr   string `xml:"listen_addr"`
}

type xmlLevelDB struct {
	CacheSize       int `xml:"cache_size"`
	BlockSize       int `xml:"block_size"`
	WriteBufferSize int `xml:"write_buffer_size"`
	MaxOpenFiles    int `xml:"max_open_files"`
}
