package ravenChecks

// ping check command
import (
  "log"
  "fmt"
  "regexp"
  "strconv"
  ."../ravenTypes"
)

var reFping []*regexp.Regexp

func init() {
  if CheckFunc == nil {
    CheckFunc = make( CheckFMap)
    CheckInit = make( CheckIMap)
  }
  CheckInit["fping"] = FpingInit
  CheckFunc["fping"] = Fping
  r, _ := regexp.Compile(`min/avg/max = (\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  reFping = append( reFping, r)
  r, _ = regexp.Compile(`xmt/rcv/\%loss = (\d+)/(\d+)/(\d+)\%,`)
  reFping = append( reFping, r)
}

func FpingInit( kw Kwargs) interface{} {
  log.Printf( "Init: %v", kw)
  return new(interface{})
}

func Fping( he *HostEntry, opts interface{}) *ExitReturn {
  e:=new(ExitReturn)

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  rtnExit, output := runExternal( "/usr/bin/fping", "-c", "8", target)

  switch rtnExit {
  case 0:
    rtt := reFping[0].FindAllStringSubmatch(output, -1)
    pls := reFping[1].FindAllStringSubmatch(output, -1)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
    loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
    e.Exit = 0
    e.Text = "Fping Okay"
    e.Perf = fmt.Sprintf( "RTT Average: %f, Loss: %d", rttAvg, loss)
    e.Long = ""
  default:
    e.Exit = 3
    e.Text = "Fping Unknown"
    e.Perf = ""
    e.Long = ""
  }

  log.Printf( "%s(Fping) exit:%d out=%s, perf=%s", he.Hostname,
    e.Exit, e.Text, e.Perf)
  return e
}
