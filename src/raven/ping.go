package raven

// ping check command

import (
  "log"
  "fmt"
  "regexp"
  "strconv"
)

var rePing []*regexp.Regexp

func init() {
  r, _ := regexp.Compile(`(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  rePing = append( rePing, r)
  r, _ = regexp.Compile(`(\d+)\% packet loss`)
  rePing = append( rePing, r)
}

func Ping( he HostEntry, opts map[string]string) (int, [3]string) {
  var rtnOut  [3]string // 0 = text; 1 = perf; 2 = extended text

  rtnExit, output := runExternal( "/usr/bin/ping", "-c", "5", he.Hostname)

  switch rtnExit {
  case 0:
    rtt := rePing[0].FindAllStringSubmatch(output, -1)
    pls := rePing[1].FindAllStringSubmatch(output, -1)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
    loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
    rtnOut[0] = "Ping Okay"
    rtnOut[1] = fmt.Sprintf( "RTT Average: %f, Loss: %d", rttAvg, loss)
    rtnOut[2] = ""
  default:
    rtnExit = 3
    rtnOut[0] = "Ping Unknown"
    rtnOut[1] = ""
    rtnOut[2] = ""
  }

  log.Printf( "%s(Ping) exit:%d out=%s, perf=%s", he.Hostname,
    rtnExit, rtnOut[0], rtnOut[1])
  return rtnExit, rtnOut
}

