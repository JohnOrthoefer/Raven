package ravenChecks

// ping check command
import (
  "log"
  "fmt"
  "regexp"
  "strconv"
  "../ravenTypes"
)

var reFping []*regexp.Regexp

func init() {
  if CheckFunc == nil {
    CheckFunc = make( map[string]func( ravenTypes.HostEntry, map[string]string) (int, [3]string))
  }
  CheckFunc["fping"] = Fping
  r, _ := regexp.Compile(`min/avg/max = (\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  reFping = append( reFping, r)
  r, _ = regexp.Compile(`xmt/rcv/\%loss = (\d+)/(\d+)/(\d+)\%,`)
  reFping = append( reFping, r)
}

func Fping( he ravenTypes.HostEntry, opts map[string]string) (int, [3]string) {
  var rtnOut  [3]string // 0 = text; 1 = perf; 2 = extended text

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  rtnExit, output := runExternal( "/usr/bin/fping", "-c", "8", target)

  //log.Printf( "Exit:%d Output:%s", rtnExit, output)
  switch rtnExit {
  case 0:
    rtt := reFping[0].FindAllStringSubmatch(output, -1)
    pls := reFping[1].FindAllStringSubmatch(output, -1)
    //log.Printf( "rtt:%v", rtt)
    //log.Printf( "pls:%v", pls)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
    loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
    rtnOut[0] = "Fping Okay"
    rtnOut[1] = fmt.Sprintf( "RTT Average: %f, Loss: %d", rttAvg, loss)
    rtnOut[2] = ""
  default:
    rtnExit = 3
    rtnOut[0] = "Fping Unknown"
    rtnOut[1] = ""
    rtnOut[2] = ""
  }

  log.Printf( "%s(Fping) exit:%d out=%s, perf=%s", he.Hostname,
    rtnExit, rtnOut[0], rtnOut[1])
  return rtnExit, rtnOut
}

