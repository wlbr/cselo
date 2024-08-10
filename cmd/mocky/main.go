package main

import (
	"bufio"
	"bytes"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alitto/pond"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

var (
	//Version is a linker injected variable for a git revision info used as version info
	Version = "Unknown build"
	//BuildTimestamp is a linker injected variable for a buildtime timestamp used in version info
	BuildTimestamp = "unknown build timestamp."

	config  *elo.Config
	maxconn = 50
)

func sendLine(line string) {
	posturl := "http://localhost:42820"
	log.Info("Posting '%s' to %s", line, posturl)
	r, err := http.NewRequest("POST", posturl, bytes.NewBufferString(line))
	if err != nil {
		log.Error("Problem creqating request for logline '%s':  %s", line, err)
		return
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          maxconn * 2,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       maxconn,
			ForceAttemptHTTP2:     false,
		}}
	res, err := client.Do(r)
	defer res.Body.Close()

	if err != nil {
		log.Error("Problem posting logline to server. '%s':  %s", line, err)
		return
	}
	if res.StatusCode != 200 {
		log.Error("Response status code is:  %d-%s", res.StatusCode, res.Status)
		return
	}
}

func ClientLoop(filename string, sendInParallel bool, pool *pond.WorkerPool) {
	log.Info("Starting mock sender loop.")
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening import CsLogFile '%s':  %s", filename, err)
	}

	scanner := bufio.NewScanner(f)
	lineno := 0

	for scanner.Scan() {
		lineno++
		buf := scanner.Text()
		if sendInParallel {
			pool.Submit(func() { sendLine(buf) })
		} else {
			sendLine(buf)
		}
	}
}

func main() {
	var sendInParallel bool
	flag.BoolVar(&sendInParallel, "parallel", true, "try to send the requests in parallel (pooled).")
	config = new(elo.Config)
	config.Initialize(Version, BuildTimestamp)
	defer config.CleanUp()

	log.Warn("Starting up")
	start := time.Now()
	pool := pond.New(maxconn, 1000)

	if config.Elo.ImportFileName == "" {
		log.Warn("No file to send to server given. Quitting...")
	} else {
		ClientLoop(config.Elo.ImportFileName, sendInParallel, pool)
	}

	if sendInParallel {
		pool.StopAndWait()
	}
	// for {
	// 	waitingtasks := pool.WaitingTasks()
	// 	activeworkers := pool.RunningWorkers()
	// 	idleworkers := pool.IdleWorkers()
	// 	fmt.Printf("%d - %d- %d\n ", activeworkers, idleworkers, waitingtasks)

	// 	if activeworkers-idleworkers == 0 && waitingtasks == 0 {
	// 		break
	// 	} else {
	// 		time.Sleep(500 * time.Millisecond)
	// 	}

	// }
	end := time.Now()
	elapsed := end.Sub(start)
	log.Warn("Shutting down, been up for %s\n", elapsed)
	//}
}
