package ravenChecks

// ping check command

import (
  "log"
  "fmt"
  "regexp"
  "strconv"
  ."../ravenTypes"
)

var rePing []*regexp.Regexp

func init() {
  if CheckFunc == nil {
    CheckFunc = make( CheckFMap)
    CheckInit = make( CheckIMap)
  }
  CheckInit["ping"] = PingInit
  CheckFunc["ping"] = Ping

  r, _ := regexp.Compile(`(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  rePing = append( rePing, r)
  r, _ = regexp.Compile(`(\d+)\% packet loss`)
  rePing = append( rePing, r)
}

func PingInit( kw Kwargs) interface{} {
  log.Printf( "Init: %v", kw)
  return new(interface{})
}

func Ping( he *HostEntry, opts interface{}) *ExitReturn {
  e:=new(ExitReturn)

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  rtnExit, output := runExternal( "/usr/bin/ping", "-c", "8", target)

  switch rtnExit {
  case 0:
    rtt := rePing[0].FindAllStringSubmatch(output, -1)
    pls := rePing[1].FindAllStringSubmatch(output, -1)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
    loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
    e.Exit = 0
    e.Text = "Ping Okay"
    e.Perf = fmt.Sprintf( "RTT Average: %f, Loss: %d", rttAvg, loss)
    e.Long = ""
  default:
    e.Exit = 3
    e.Text = "Ping Unknown"
    e.Perf = ""
    e.Long = ""
  }

  log.Printf( "%s(Ping) exit:%d out=%s, perf=%s", he.Hostname,
    e.Exit, e.Text, e.Perf)
  return e
}

