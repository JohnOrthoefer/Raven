package raven

import (
  "time"
  "log"
  "regexp"
  ."./ravenTypes"
  ."./ravenChecks"
)

type hostMap map[string]*HostEntry
type checkMap map[string]*CheckEntry
var hosts hostMap
var checks checkMap

func init() {
  hosts = make( hostMap)
  checks = make( checkMap)
}

func contains( v string, s []string) bool {
  for _,i := range s {
    if i == v {
      return true
    }
  }
  return false
}

func isHost( h string) bool {
  _,ok := hosts[h]
  return ok
}

func getEntry( kv Kwargs, n string) string {
  if v,ok := kv[n]; ok {
    return v
  }
  return ""
}

func newHost( n string, kv Kwargs) *HostEntry {
  r := new( HostEntry)
  r.DisplayName = n
  r.IPv4 = getEntry( kv, "ipv4")
  r.Hostname = getEntry( kv, "hostname")
  r.Group = getEntry( kv, "group")
  return r
}

func newCheck( n string, kv Kwargs) *CheckEntry {
  r := new( CheckEntry)
  r.DisplayName = n
  // Check function that will be run
  r.CheckN = getEntry( kv, "checkwith")
  r.CheckF = CheckFunc[r.CheckN]

  // set up the run intervals
  t,_ := time.ParseDuration( "30s")
  for i:= range r.Interval {
    r.Interval[i] = t
  }
  re := regexp.MustCompile( `\s+`)
  k:=getEntry(kv, "interval")
  inter := re.Split( k, -1)
  for i,j := range inter {
    if t,ok := time.ParseDuration( j); ok==nil {
      r.Interval[i] = t
    } else {
      log.Printf( "Error Parsing %s", j)
    }
  }

  // array of hosts that use this check
  for _,n := range re.Split( getEntry(kv, "hosts"), -1) {
    r.Hosts = append(r.Hosts, n)
  }

  // move anything else random (which will be used by the check command 
  // into basically a kwargs structure
  Options := make( Kwargs)
  for k,v := range kv {
    if !contains( k, []string{"checkwith", "interval", "hosts"}) {
      Options[k] = v
    }
  }

  r.Options = CheckInit[r.CheckN]( Options)
  return r
}

func AddEntry( n string, kv map[string]string) {
  if _,ok := kv["hostname"]; ok {
    hosts[n] = newHost(n, kv)
  } else if _,ok := kv["checkwith"]; ok {
    checks[n] = newCheck(n, kv)
  } else {
    log.Printf( "Unknown Section Type %s", n)
  }
}

func GetCheckEntry( c string) *CheckEntry {
  if _,ok := checks[c]; !ok {
    return nil
  }
  return checks[c]
}

func ListChecks() []string {
  rtn:=[]string{}
  for n:=range checks {
    rtn = append(rtn, n)
  }
  return rtn
}

func ListCheckHosts( c string) []string {
  if _,ok := checks[c]; ok {
    var rtn []string
    for _,h:=range checks[c].Hosts {
      if !isHost(h) {
        log.Printf( "%s, can not find %s", c, h)
        continue
      }
      rtn = append( rtn, h)
    }
    return rtn
  }
  return nil
}

func GetHostEntry( c string) *HostEntry {
  if _,ok := hosts[c]; ok {
    return hosts[c]
  }
  return nil
}


func printHosts() {
  for i:= range hosts {
    log.Printf( "hosts[%s] = %v", i, hosts[i])
  }
}

func printChecks() {
  for i:= range checks {
    log.Printf( "checks[%s] = %v", i, checks[i])
  }
}

func DumpStorage() {
  printHosts()
  printChecks()
}
