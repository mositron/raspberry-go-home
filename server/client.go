package server

import (
	"bufio"
	"fmt"
	"github.com/positronth/raspberry-go-home/config"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func LoadClient(name string, ip string, port int) {
	//log.Println("[CLIENT] - LoadClient", name, ip, port)
	_exists := false
	for k, v := range Con {
		if v.Ip == ip {
			if v.Name == "" {
				Con[k].Name = name
			}
			_exists = true
			break
		}
	}
	if !_exists {
		go func() {
			create(ip, port)
		}()
		//Con = append(Con, Connects{ip: _ip, port: _port, found: _found, last: _found, conn: client})
	}
}
func reconnect(ip string, port int) {
	for k, v := range Con {
		if v.Ip == ip {
			Con = append(Con[:k], Con[k+1:]...)
			break
		}
	}
	log.Println("[ERROR] disconnected!. trying to connect", ip, "in 1 minute")
	time.Sleep(time.Minute)
	go func() {
		create(ip, port)
	}()
}
func create(ip string, port int) {
	//log.Println("[CLIENT] Connecting to ["+name+"] ", fmt.Sprint(ip, ":", port))
	conn, err := net.Dial("tcp", fmt.Sprint(ip, ":", port))
	if err != nil {
		//log.Println("[ERROR] Can not connect to ["+name+"] ", fmt.Sprint(ip, ":", port))
		return
	}
	m := []string{config.Host, string(config.PortWeb), config.Location, config.Cam1, config.Cam2, config.Cam3, config.Sound1, config.Sound2, config.IRRev, config.IRSend}
	conn.Write([]byte("INIT " + strings.Join(m, "###") + "\n"))
	log.Println("[CLIENT]", fmt.Sprint(ip, ":", port), "is connected")
	Con = append(Con, &Connects{Ip: ip, Port: port, Found: time.Now(), Last: time.Now(), Conn: conn, IsClient: true})
	//isClient
	go func() {
		for {
			time.Sleep(time.Minute)
			_, err := conn.Write([]byte("PING " + strconv.FormatInt(time.Now().Unix(), 10) + "\n"))
			if err != nil {
				if err == io.EOF {
					go func() {
						reconnect(ip, port)
					}()
					//} else {
				}
				log.Println("[ERROR] client.go > create > Write > ping :", err.Error())
				break
			}
		}
	}()
	for {
		ms, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				go func() {
					reconnect(ip, port)
				}()
			}
			log.Println("[ERROR] client.go > create > readstring :", err.Error())
			break
		} else {
			packet(conn, strings.TrimSpace(ms))
		}
	}
}
