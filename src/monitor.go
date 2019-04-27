package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type hostRecord struct {
	hostname string
	last     time.Time
	exitCode int
	inflight bool
}

var msgs = make(chan *hostRecord)
var done = make(chan bool)
var list []*hostRecord

// find the oldest record that is due to be run
func getDue() *hostRecord {
	dur, _ := time.ParseDuration("30s")
	var rtn *hostRecord
	rtn = nil

	for i := range list {
		if list[i].inflight {
			continue
		}
		if time.Since(list[i].last) < dur {
			continue
		}
		if rtn == nil {
			rtn = list[i]
			continue
		}
		if rtn.last.After(list[i].last) {
			rtn = list[i]
		}
	}
	return rtn
}

func getTimestamp() string {
	return time.Now().Format(time.RFC1123Z)
}

func LogMessage(t string) {
	fmt.Printf("%s %s\n", getTimestamp(), t)
}

func produce() {
	for {
		i := getDue()
		if i == nil {
			LogMessage("Sleeping...")
			time.Sleep(time.Second * 5)
		} else {
			LogMessage(fmt.Sprintf("%s last was %s", i.hostname, i.last.Format(time.RFC1123Z)))
			i.inflight = true
			msgs <- i
		}
	}
}

func consume(work int) {
	r, _ := regexp.Compile(`(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
	st, _ := regexp.Compile(`(\d+)\% packet loss`)
	for {
		msg := <-msgs
		LogMessage(fmt.Sprintf("Worker %d, pinging %s", work, msg.hostname))
		cmd := exec.Command("/usr/bin/ping", "-c", "5", msg.hostname)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		msg.exitCode = 0
		if err == nil {
			msg.exitCode = 3
		}
		rtt := r.FindAllStringSubmatch(out.String(), -1)
		pls := st.FindAllStringSubmatch(out.String(), -1)
		rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
		loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
		LogMessage(fmt.Sprintf("Worker %d: Done %s rtt:%f loss:%d", work, msg.hostname, rttAvg, loss))
		msg.last = time.Now()
		msg.inflight = false
	}
}

func main() {
	h := []string{"127.0.0.1", "www.google.com", "www.cnn.com", "www.disney.com", "172.17.2.254"}

	for i := range h {
		t := new(hostRecord)
		t.hostname = h[i]
		t.last = time.Unix(0, 0)
		list = append(list, t)
	}
	//fmt.Printf( "%q\n",list)
	go produce()
	go consume(1)
	go consume(2)
	<-done
}