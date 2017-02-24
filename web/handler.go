package web

import (
	"github.com/positronth/raspberry-go-home/config"
	"github.com/positronth/raspberry-go-home/server"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

type MyHandler struct {
	sync.Mutex
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", config.Name)
	cmd := strings.Split(strings.ToLower(strings.TrimLeft(r.URL.Path, "/")), "/")[0]
	if cmd == "files" {
		dat, err := ioutil.ReadFile(r.URL.Path[1:])
		if err == nil {
			ext := filepath.Ext(r.URL.Path)
			if ext == ".css" {
				w.Header().Add("Content-Type", "text/css; charset=utf-8")
			} else if ext == ".js" {
				w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
			}
			w.Write([]byte(dat))
		}
	} else {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		ct := "<h3 class='bar'><i class='fa fa-sitemap'></i> IPs - LAN</h3>"
		ct += "<div class='table-responsive'><table class='table table-bordered table-condensed' width='100%'><thead><tr><th>Name</th><th>IP</th><th class='time'>Added</th><th class='time'>Updated</th></tr><thead><tbody>"
		for _, v := range server.Lan {
			ct += "<tr><td><i class='fa fa-television'></i> " + v.Name + "</td><td>" + v.Ip + "</td>"
			ct += "<td class='time'>" + v.Found.Format("2006-01-02 15:04:05") + "</td>"
			ct += "<td class='time'>" + v.Last.Format("2006-01-02 15:04:05") + "</td>"
			ct += "</tr>"
		}
		ct += "</tbody></table></div>"
		ct += "<h3 class='bar'><i class='fa fa-exchange'></i> Connections</h3>"
		ct += "<div class='table-responsive'><table class='table table-bordered table-condensed' width='100%'><thead><tr><th>Name</th><th>IP</th><th>Location</th><th class='i'><i class='fa fa-camera'></i></th><th class='i'><i class='fa fa-volume-up'></i></th><th class='i'><i class='fa fa-bolt'></i></th><th class='time'>Added</th><th class='time'>Updated</th></tr><thead><tbody>"
		for _, v := range server.Con {
			ct += "<tr><td><i class='fa fa-wifi'></i> " + v.Name + "</td><td>" + v.Ip + "</td>"
			ct += "<td>" + v.Location + "</td>"
			ct += "<td class='i'>-</td><td class='i'>-</td><td class='i'>-</td>"
			ct += "<td class='time'>" + v.Found.Format("2006-01-02 15:04:05") + "</td>"
			ct += "<td class='time'>" + v.Last.Format("2006-01-02 15:04:05") + "</td>"
			ct += "</tr>"
		}
		ct += "</tbody></table></div>"
		dat, err := ioutil.ReadFile("files/html/index.html")

		if err == nil {
			html := strings.Replace(string(dat), "{CONTENT}", ct, -1)
			html = strings.Replace(html, "{TITLE}", config.Name, -1)
			w.Write([]byte(html))
		} else {
			w.Write([]byte(err.Error()))
		}
	}
}
