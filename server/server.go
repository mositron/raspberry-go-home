package server

import (
	"bufio"
	"fmt"
	"home/config"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	count int
}

func (s *Server) Listen(port int) {
	s.count++
	Con = []*Connects{}
	log.Println("[SERVER] Launching SERVER on port", fmt.Sprint(":", port))
	go func() {
		for {
			time.Sleep(time.Minute)
			log.Println("##### Connections #####")
			for _, v := range Con {
				log.Println(v.Name, v.Ip, v.Port, v.IsClient)
			}
			log.Println("#######################")
		}
	}()

	// listen on all interfaces
	serv, _ := net.Listen("tcp", fmt.Sprint(":", port))
	conns := s.clientConns(serv)
	for {
		go s.handleConn(<-conns)
	}
}

func (s *Server) clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				log.Printf("[SERVER] couldn't accept:" + err.Error())
				continue
			}

			m := []string{config.Host, string(config.PortWeb), config.Location, config.Cam1, config.Cam2, config.Cam3, config.Sound1, config.Sound2, config.IRRev, config.IRSend}
			client.Write([]byte("INIT " + strings.Join(m, "###") + "\n"))
			time.Sleep(time.Second)
			_ip := client.RemoteAddr().(*net.TCPAddr).IP.String()
			_port := client.RemoteAddr().(*net.TCPAddr).Port
			_found := time.Now()
			_exists := false
			for key, val := range Con {
				if val.Ip == _ip {
					Con[key].Port = _port
					Con[key].Last = _found
					Con[key].Conn = client
					_exists = true
					break
				}
			}
			if !_exists {
				Con = append(Con, &Connects{Ip: _ip, Port: _port, Found: _found, Last: _found, Conn: client, IsClient: false})
			}
			ch <- client
		}
	}()
	return ch
}

func (s *Server) handleConn(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadString('\n')
		if err != nil { // EOF, or worse
			break
		}
		_ip := client.RemoteAddr().(*net.TCPAddr).IP.String()
		_found := time.Now()
		for key, val := range Con {
			if val.Ip == _ip {
				Con[key].Last = _found
				break
			}
		}
		packet(client, strings.TrimSpace(line))
	}
}
