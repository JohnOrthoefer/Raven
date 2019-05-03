package raven

import (
  "time"
  "fmt"
  "regexp"
  "strings"
  "sort"
  "./ravenLog"
  "./ravenTypes"
  "./ravenChecks"
)

type hostMap map[string]*ravenTypes.HostEntry
type checkMap map[string]*ravenTypes.CheckEntry
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

func newHost( n string, kv ravenTypes.Kwargs) *ravenTypes.HostEntry {
  if !kv.GetKwargBool("enabled",true) {
    ravenLog.SendMessage( 10, "newHost",
      fmt.Sprint( "Host %s is disabled", n))
    return nil
  }
  r := new( ravenTypes.HostEntry)
  r.DisplayName = n
  r.IPv4 = kv.GetKwargStrTrim( "ipv4", "")
  r.Hostname = getEntry( kv, "hostname", true)
  r.Group = getEntry( kv, "group", true)
  return r
}

func newCheck( n string, kv ravenTypes.Kwargs) *ravenTypes.CheckEntry {
  r := new( ravenTypes.CheckEntry)
  re := regexp.MustCompile( `\s+`)

  r.DisplayName = n
  // Check function that will be run
  r.CheckN = getEntry( kv, "checkwith", true)
  r.CheckF = ravenChecks.CheckFunc[r.CheckN]

  // set up the run intervals
  t,_ := time.ParseDuration( "30s")
  for i:= range r.Interval {
    r.Interval[i] = t
  }
  k:=getEntry(kv, "interval", true)
  inter := re.Split( k, len(r.Interval))
  for i,j := range inter {
    ravenLog.SendError( 99, "newCheck", fmt.Sprintf( "Parsing %s", j))
    if t,ok := time.ParseDuration( j); ok==nil {
      r.Interval[i] = t
    } else {
      ravenLog.SendError( 10, "newCheck", fmt.Sprintf( "Error Parsing %s", j))
    }
  }

  r.Threshold = int(kv.GetKwargInt( "threshold", 5))

  // array of hosts that use this check
  for _,n := range re.Split( getEntry(kv, "hosts", true), -1) {
    // dedup the hosts
    if contains( n, r.Hosts) {
      continue
    }
    r.Hosts = append(r.Hosts, n)
  }

  // move anything else random (which will be used by the check command 
  // into basically a kwargs structure
  Options := make( ravenTypes.Kwargs)
  for k,v := range kv {
    k = strings.ToLower(k)
    if !contains( k, []string{"checkwith", "interval", "hosts", "threshold"}) {
      Options[k] = v
    }
  }

  r.Options = ravenChecks.CheckInit[r.CheckN]( Options)
  return r
}

func AddEntry( n string, kv map[string]string) {
  if _,ok := kv["hostname"]; ok {
    hosts[n] = newHost(n, kv)
  } else if _,ok := kv["checkwith"]; ok {
    checks[n] = newCheck(n, kv)
  } else {
    ravenLog.SendMessage( 10, "AddEntry", fmt.Sprintf( "Unknown Section Type %s", n))
  }
}

func GetCheckEntry( c string) *ravenTypes.CheckEntry {
  if _,ok := checks[c]; !ok {
    return nil
  }
  return checks[c]
}

func ListChecks() []string {
  rtn:=sort.StringSlice{}
  for n:=range checks {
    rtn = append(rtn, n)
  }
  rtn.Sort()
  return rtn
}

func ListCheckHosts( c string) []string {
  if _,ok := checks[c]; ok {
    rtn:=sort.StringSlice{}
    for _,h:=range checks[c].Hosts {
      if !isHost(h) {
        ravenLog.SendError( 10, "ListCheckHosts", fmt.Sprintf( "%s, can not find %s", c, h))
        continue
      }
      rtn = append( rtn, h)
    }
    rtn.Sort()
    return rtn
  }
  return nil
}

func GetHostEntry( c string) *ravenTypes.HostEntry {
  if _,ok := hosts[c]; ok {
    return hosts[c]
  }
  return nil
}


func printHosts() {
  for i:= range hosts {
    ravenLog.SendMessage( 10, "printHosts", fmt.Sprintf( "hosts[%s] = %v", i, hosts[i]))
  }
}

func printChecks() {
  for i:= range checks {
    ravenLog.SendMessage( 10, "printChecks", fmt.Sprintf( "checks[%s] = %v", i, checks[i]))
  }
}

func DumpStorage() {
  printHosts()
  printChecks()
}
