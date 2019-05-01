package ravenChecks

// ping check command

import (
  "fmt"
  "regexp"
  "strconv"
  "../ravenLog"
  "../ravenTypes"
)

var rePing []*regexp.Regexp

type pingOpts struct {
  pingProg  string
  rttWarn   float64
  rttCrit   float64
  lossWarn  int64
  lossCrit  int64
  count     string
}

func init() {
  registerHandler( "ping", PingInit, Ping)

  r, _ := regexp.Compile(`(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)/(\d+\.?\d+)`)
  rePing = append( rePing, r)
  r, _ = regexp.Compile(`(\d+)\% packet loss`)
  rePing = append( rePing, r)
}

func PingInit( kw ravenTypes.Kwargs) interface{} {
  var r interface{}
  rtn := new( pingOpts)
  rtn.pingProg = kw.GetKwargStr( "program", "/usr/bin/ping")
  rtn.rttWarn = kw.GetKwargFloat( "rtt_warn", 20.0)
  rtn.lossWarn = kw.GetKwargInt( "loss_warn", 20)
  rtn.rttCrit = kw.GetKwargFloat( "rtt_crit", 30.0)
  rtn.lossCrit = kw.GetKwargInt( "loss_crit", 40)
  rtn.count = kw.GetKwargStr( "count", "5")
  r = rtn
  return r
}

func Ping( he *ravenTypes.HostEntry, options interface{}) *ravenTypes.ExitReturn {
  e:=new(ravenTypes.ExitReturn)
  opts := options.(*pingOpts)

  target := he.Hostname
  if he.IPv4 != "" {
    target = he.IPv4
  }
  rtnExit, output := runExternal( opts.pingProg, "-c", opts.count, target)

  switch rtnExit {
  case 0:
    rtt := rePing[0].FindAllStringSubmatch(output, -1)
    pls := rePing[1].FindAllStringSubmatch(output, -1)
    rttAvg, _ := strconv.ParseFloat(rtt[0][3], 32)
    loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
    e.Exit = 0
    e.Text = "Ping Okay"
    e.Perf = fmt.Sprintf( "RTT Average: %4.2f, Loss: %d", rttAvg, loss)
    e.Long = fmt.Sprintf( "Count:%s Warn:%d,%4.2f Crit:%d,%4.2f",
              opts.count, opts.lossWarn, opts.rttWarn,
              opts.lossCrit, opts.rttCrit)
    if opts.rttWarn < rttAvg || opts.lossWarn < loss {
      e.Exit = 1
      e.Text = "Ping Warning"
      e.Long = fmt.Sprintf( "WARNING %4.2f < %4.2f or %d < %d",
        opts.rttWarn, rttAvg, opts.lossWarn, loss)
    }
    if opts.rttCrit < rttAvg || opts.lossCrit < loss {
      e.Exit = 2
      e.Text = "Ping Critical"
      e.Long = fmt.Sprintf( "CRITICAL %4.2f < %4.2f or %d < %d",
        opts.rttCrit, rttAvg, opts.lossCrit, loss)
    }
  default:
    e.Exit = 3
    e.Text = "Ping Unknown"
    e.Perf = ""
    e.Long = output
  }

  ravenLog.SendMessage( 10, "ping", fmt.Sprintf( "%s(Ping) exit:%d out=%s, perf=%s", he.Hostname,
    e.Exit, e.Text, e.Perf))
  return e
}

