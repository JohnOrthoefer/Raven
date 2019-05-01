package ravenChecks

// Nagios externals check command
import (
  "fmt"
  //"regexp"
  "strings"
  "../ravenTypes"
  "../ravenLog"
)

type nagiosOpts struct {
  prog      string
  progOpts  []string
  //resplit   *regexp.Regexp
}

func init() {
  registerHandler( "nagios", nagiosInit, nagios)
}

func nagiosInit( kw ravenTypes.Kwargs) interface{} {
  var r interface{}
  rtn := new( nagiosOpts)
  rtn.prog = kw.GetKwargStr( "program", "/usr/lib/monitoring-plugins/check_ping")
  rtn.progOpts = []string{ "-w", "20,20%", "-c", "40,40%" }
  //rtn.resplit = regexp.MustCompile("|")
  r = rtn
  return r
}

func nagios( he *ravenTypes.HostEntry, options interface{}) *ravenTypes.ExitReturn {
  e:=new(ravenTypes.ExitReturn)
  opts := options.(*nagiosOpts)

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  fullOpts := append( opts.progOpts, "-H", target)
  rtnExit, output := runExternal( opts.prog, fullOpts...)

  switch rtnExit {
  case 0:
    fallthrough
  case 1:
    fallthrough
  case 2:
    e.Exit = rtnExit
    s := strings.Split( output, "|")
    switch len(s) {
      case 3:
        e.Long = s[2]
        fallthrough
      case 2:
        e.Perf = s[1]
        fallthrough
      case 1:
        e.Text = s[0]
      default:
        ravenLog.SendError(10, "Nagios", "Failed to split output")
    }
  default:
    e.Exit = 3
    e.Text = "Nagios Unknown"
    e.Perf = ""
    e.Long = output
  }

  ravenLog.SendMessage(10, "Nagios", fmt.Sprintf( "%s(Nagios) exit:%d out=%s, perf=%s, long=%s", he.Hostname,
    e.Exit, e.Text, e.Perf, e.Long))
  return e
}
