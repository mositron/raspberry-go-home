package server

import (
	"errors"
	"home/config"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func DynamicIP() {
	resp, err := http.Get("http://mokoidea.com/pi/ip.php")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func IpsLan() {
	hosts, _ := Hosts(config.Ip + "/24")
	concurrentMax := 100
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan Pong, len(hosts))
	doneChan := make(chan []Pong)

	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
		//fmt.Println("sent: " + ip)
	}
	alives := <-doneChan

	//	Lan = []*Connects{}
	now := time.Now()
	for _, v := range alives {
		found := false
		for k2, v2 := range Lan {
			if v.Ip == v2.Ip {
				found = true
				Lan[k2].Last = now
				break
			}
		}
		if !found {
			Lan = append(Lan, &Connects{Ip: v.Ip, Name: v.Name, Found: now, Last: now})
		}
	}
	for i := len(Lan) - 1; i >= 0; i-- {
		if Lan[i].Last != now {
			Lan = append(Lan[:i], Lan[i+1:]...)
		}
	}
	log.Println("###### IPs - Lan ######")
	for _, v := range Lan {
		log.Println(v.Ip, v.Name)
	}
	log.Println("#######################")
	for _, v := range Lan {
		//	if cf.ip != v.Ip {
		LoadClient(v.Name, v.Ip, config.PortServer)
		//	}
	}
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip = ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type Pong struct {
	Ip    string
	Alive bool
	Name  string
	Txt   string
}

func ping(pingChan <-chan string, pongChan chan<- Pong) {
	for ip := range pingChan {
		//txt, err := exec.Command("ping", "-c1", "-t1", ip).Output()
		cmd := ""
		isLinux := true
		if runtime.GOOS == "windows" {
			cmd = "-a -n 1"
			isLinux = false
		} else {
			cmd = "-c1 -t1"
		}
		txt, err := exec.Command("ping", cmd, ip).Output()
		//log.Println("ping", cmd, ip)
		var alive bool
		name := ""
		t := CToGoString(txt)
		if err != nil {
			alive = false
			//log.Println("[ERROR] ping", cmd, ip, err.Error(), t)
		} else {
			reply := "Reply from " + ip + ": bytes="
			if isLinux {
				reply = "bytes from " + ip + ":"
			}
			//log.Println(ip, t)
			if strings.Contains(t, reply) {
				alive = true
				addr, _ := net.LookupAddr(ip)
				name = strings.Join(addr, " ")
				if !isLinux {
					ln := strings.Split(t, "\n")
					for _, v := range ln {
						if strings.Contains(v, "Pinging ") {
							v = strings.Trim(v, " ")
							n := strings.Split(v, " ")
							if strings.Contains(n[2], "["+ip+"]") {
								n = strings.Split(n[1], ".")
								name = n[0]
							}
						}
					}
				}
			} else {
				alive = false
			}
		}
		pongChan <- Pong{Ip: ip, Alive: alive, Name: name, Txt: t}
	}
}

func receivePong(pongNum int, pongChan <-chan Pong, doneChan chan<- []Pong) {
	var alives []Pong
	for i := 0; i < pongNum; i++ {
		pong := <-pongChan
		//fmt.Println("received:", pong)
		if pong.Alive {
			alives = append(alives, pong)
		}
	}
	doneChan <- alives
}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
