package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	Name       string
	Host       string
	Location   string
	Cam1       string
	Cam2       string
	Cam3       string
	Sound1     string
	Sound2     string
	IRRev      string
	IRSend     string
	PortWeb    int
	PortServer int
	PortCam    int
	Ip         string
	//is_master   bool
)

func LoadConf(f string) {
	f2, _ := ioutil.ReadFile(f)
	var dat map[string]interface{}
	if err := json.Unmarshal(f2, &dat); err != nil {
		panic(err)
	}
	host, err := os.Hostname()
	if err == nil {
		Host = host
	}
	PortWeb = int(dat["port_web"].(float64))
	PortServer = int(dat["port_server"].(float64))
	PortCam = int(dat["port_cam"].(float64))
	Location = dat["location"].(string)
	Cam1 = dat["cam1"].(string)
	Cam2 = dat["cam2"].(string)
	Cam3 = dat["cam3"].(string)
	Sound1 = dat["sound1"].(string)
	Sound2 = dat["sound2"].(string)
	IRRev = dat["ir_rev"].(string)
	IRSend = dat["ir_send"].(string)
}
