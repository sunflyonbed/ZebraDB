package main

import (
	"flag"
	"fmt"
	"runtime"

	l4g "log4go"

	"common"
)

var (
	gConf = new(xmlConfig)
	gDB   *LevelDB
)

var configFile = flag.String("config", "../config/zebra_config.xml", "")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if err := common.LoadConfig(*configFile, gConf); err != nil {
		panic(fmt.Sprintf("load config %v fail: %v", *configFile, err))
	}
	l4g.LoadConfiguration(gConf.Log.Config)
	defer l4g.Close()

	l4g.Debug("configure: %+v", *gConf)

	var err error
	gDB, err = NewLevelDB(gConf.Server.Path)
	if err != nil {
		panic(err.Error())
	}
	defer gDB.Close()
	go RedisServer(gConf.Server.ListenAddr)
	l4g.Info("monitor redis info: %s %d", gConf.Server.MonitorAddr, gConf.Server.MonitorIndex)
	RedisMonitor(gConf.Server.MonitorAddr, gConf.Server.MonitorIndex)
}
