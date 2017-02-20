package server

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var Con []*Connects
var Lan []*Connects

type Connects struct {
	Name string
	Ip   string
	Port int

	Web      int
	Location string
	Cam1     string
	Cam2     string
	Cam3     string
	Sound1   string
	Sound2   string
	IRRev    string
	IRSend   string

	Found    time.Time
	Last     time.Time
	Conn     net.Conn
	IsClient bool
}

func packet(c net.Conn, t string) {
	ip := c.RemoteAddr().(*net.TCPAddr).IP.String()
	key := -1
	name := ""
	for k, v := range Con {
		if v.Ip == ip {
			Con[k].Last = time.Now()
			if Con[k].Name != "" {
				name = "[" + Con[k].Name + "]"
			}
			key = k
			break
		}
	}
	if name == "" {
		name = ip
	}
	s := strings.SplitN(t, " ", 2)
	cmd, txt := s[0], s[1]
	log.Println("[MESSAGE] " + name + " : " + t)
	switch cmd {
	case "INIT":
		if key != -1 {
			//m := []string{config.Host, string(config.PortWeb), config.Location, config.Cam1, config.Cam2, config.Cam3, config.Sound1, config.Sound2, config.IRRev, config.IRSend}
			m := strings.Split(txt, "###")
			Con[key].Name = m[0]
			Con[key].Web, _ = strconv.Atoi(m[1])
			Con[key].Location = m[2]
			Con[key].Cam1 = m[3]
			Con[key].Cam2 = m[4]
			Con[key].Cam3 = m[5]
			Con[key].Sound1 = m[6]
			Con[key].Sound2 = m[7]
			Con[key].IRRev = m[8]
			Con[key].IRSend = m[9]
		}
	case "PING":
		c.Write([]byte("PONG " + txt + "\n"))
	case "PONG":
	default:
	}
}
