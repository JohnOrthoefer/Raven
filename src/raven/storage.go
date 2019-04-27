package raven

import (
  "time"
  "log"
  "regexp"
)

type HostEntry struct {
  IPv4        string
  Hostname    string
  DisplayName string
  Group       string
}

type CheckEntry struct {
  DisplayName string
  CheckF      func( h HostEntry, p map[string]string)
  CheckN      string
  Interval    [4]time.Duration
  Hosts       []string
  Options     map[string]string
}

var hosts map[string]*HostEntry
var checks map[string]*CheckEntry

func contains( v string, s []string) bool {
  for _,i := range s {
    if i == v {
      return true
    }
  }
  return false
}

func getEntry( kv map[string]string, n string) string {
  if v,ok := kv[n]; ok {
    return v
  }
  return ""
}

func newHost( n string, kv map[string]string) *HostEntry {
  r := new( HostEntry)
  r.DisplayName = n
  r.IPv4 = getEntry( kv, "ipv4")
  r.Hostname = getEntry( kv, "hostname")
  r.Group = getEntry( kv, "group")
  return r
}

func newCheck( n string, kv map[string]string) *CheckEntry {
  r := new( CheckEntry)
  r.DisplayName = n
  // Check function that will be run
  r.CheckN = getEntry( kv, "checkwith")

  // set up the run intervals
  t,_ := time.ParseDuration( "30s")
  r.Interval[0] = t
  r.Interval[1] = t
  r.Interval[2] = t
  r.Interval[3] = t
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
  r.Options = make( map[string]string)
  for k,v := range kv {
    if !contains( k, []string{"checkwith", "interval", "hosts"}) {
      r.Options[k] = v
    }
  }
  return r
}

func AddEntry( n string, kv map[string]string) {
  if hosts == nil {
    hosts = make( map[string]*HostEntry)
  }
  if checks == nil {
    checks = make( map[string]*CheckEntry)
  }

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
    return checks[c].Hosts
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
