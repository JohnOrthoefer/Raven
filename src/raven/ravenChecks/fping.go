package ravenChecks

// ping check command
import (
  "fmt"
  "regexp"
  "strconv"
  "../ravenTypes"
  "../ravenLog"
)

var reFping []*regexp.Regexp

type fPingOpts struct {
  prog      string
  rttWarn   float64
  rttCrit   float64
  lossWarn  int64
  lossCrit  int64
  count     string
}


func init() {
  registerHandler( "fping", FpingInit, Fping)

  r, _ := regexp.Compile(`min/avg/max = (\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  reFping = append( reFping, r)
  r, _ = regexp.Compile(`xmt/rcv/\%loss = (\d+)/(\d+)/(\d+)\%,`)
  reFping = append( reFping, r)
}

func FpingInit( kw ravenTypes.Kwargs) interface{} {
  var r interface{}
  rtn := new( fPingOpts)
  rtn.prog = kw.GetKwargStr( "program", "/usr/bin/fping")
  rtn.rttWarn = kw.GetKwargFloat( "rtt_warn", 20.0)
  rtn.lossWarn = kw.GetKwargInt( "loss_warn", 20)
  rtn.rttCrit = kw.GetKwargFloat( "rtt_crit", 30.0)
  rtn.lossCrit = kw.GetKwargInt( "loss_crit", 40)
  rtn.count = kw.GetKwargStr( "count", "5")
  r = rtn
  return r
}

func Fping( he *ravenTypes.HostEntry, options interface{}) *ravenTypes.ExitReturn {
  e:=new(ravenTypes.ExitReturn)
  opts := options.(*fPingOpts)

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  rtnExit, output := runExternal( opts.prog, "-c", opts.count, target)

  switch rtnExit {
  case 0:
    rtt := reFping[0].FindAllStringSubmatch(output, -1)
    pls := reFping[1].FindAllStringSubmatch(output, -1)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 64)
    loss, _ := strconv.ParseInt(pls[0][3], 10, 64)
    e.Exit = 0
    e.Text = "Fping Okay"
    e.Perf = fmt.Sprintf( "RTT Average: %4.2f, Loss: %d", rttAvg, loss)
    e.Long = fmt.Sprintf( "Count:%s Warn:%d,%4.2f Crit:%d,%4.2f",
              opts.count, opts.lossWarn, opts.rttWarn,
              opts.lossCrit, opts.rttCrit)
    if opts.rttWarn < rttAvg || opts.lossWarn < loss {
      e.Exit = 1
      e.Text = "Fping Warning"
      e.Long = fmt.Sprintf( "WARNING %4.2f < %4.2f or %d < %d",
        opts.rttWarn, rttAvg, opts.lossWarn, loss)
    }
    if opts.rttCrit < rttAvg || opts.lossCrit < loss {
      e.Exit = 2
      e.Text = "Fping Critical"
      e.Long = fmt.Sprintf( "CRITICAL %4.2f < %4.2f or %d < %d",
        opts.rttCrit, rttAvg, opts.lossCrit, loss)
    }
  default:
    e.Exit = 3
    e.Text = "Fping Unknown"
    e.Perf = ""
    e.Long = output
  }

  ravenLog.SendMessage(10, "Check fping", fmt.Sprintf( "%s(Fping) exit:%d out=%s, perf=%s, long=%s", he.Hostname,
    e.Exit, e.Text, e.Perf, e.Long))
  return e
}
