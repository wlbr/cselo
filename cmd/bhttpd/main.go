package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/wlbr/commons/log"

	"github.com/wlbr/cselo/elo"
)

var (
	//Version is a linker injected variable for a git revision info used as version info
	Version = "Unknown build"
	//BuildTimestamp is a linker injected variable for a buildtime timestamp used in version info
	BuildTimestamp = "unknown build timestamp."
)

type httpd struct {
	config   *elo.Config
	recorder *elo.Recorder
	wg       *sync.WaitGroup
	listener net.Listener
}

func NewHttpd(cfg *elo.Config) *httpd {
	wg := new(sync.WaitGroup)
	self := &httpd{wg: wg, config: cfg, recorder: elo.NewRecorder(cfg, wg)}
	return self
}

func (h *httpd) Reply(conn net.Conn, msg string) (int, error) {
	contentlength := len(msg)
	var out string = fmt.Sprintf(`HTTP/1.1 200 OK
Server: bthttpd/1.0 (Go)
ContentContent-Length: %d
Content-Language: en
Connection: close
Content-Type: text/html

%s`, contentlength, msg)
	n, err := conn.Write([]byte(out))
	return n, err
}

func (h *httpd) Loop() {
	h.wg.Add(1)
	defer h.wg.Done()
	go h.recorder.Loop()
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			log.Error("Error accepting new connection %s", err)
		}
		buff := make([]byte, 4096)
		n, err := conn.Read(buff)
		if err != nil {
			log.Error("Failed to read from connection %s", err)
		} else {
			sbuf := string(buff[:n])
			h.recorder.Record(sbuf)
			_, err = h.Reply(conn, "processed")
			if err != nil {
				log.Error("Failed to write response %s", err)
			}
			if sbuf == "cselo:StopProcessor" {
				return
			}
		}
		conn.Close()
	}

}

func main() {

	config := new(elo.Config)
	config.Initialize(Version, BuildTimestamp)
	defer config.CleanUp()
	log.Warn("Starting up")

	httpd := NewHttpd(config)
	var err error

	httpd.listener, err = net.Listen("tcp", "127.0.0.1:"+config.Elo.Port)
	if err != nil {
		log.Fatal("Failed to start to listen: %s", err)
		config.FatalExit()
	}
	httpd.Loop()

}
