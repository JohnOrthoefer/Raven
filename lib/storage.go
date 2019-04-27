package storage

import (
  "log"
  "time"
)
/*
//type ExitRecord struct {
//  run  time.Time
//  value int
//  output string
//  perf string
//  long string
//}
*/

type CheckRecord struct {
//  check func( c CallRecord) (e ExitRecord)
  freq [4]time.Duration 
  nextRun time.Time
//  exit ExitRecord
//  opt map[string]string
}

type HostRecord struct {
  host  string
  group string
  check map[string]*CheckRecord
}

// global list of all the records
var status map[string]*HostRecord

func hostExists( dsp string) (bool) {
  _,ok := status[dsp]
  return ok
}

func NewHost( dsp string, h string, g string) {
  // make a new Map if needed
  if status == nil {
    status = make( map[string]*HostRecord)
  }

  // new means not exisiting
  if hostExists( dsp) {
    log.Printf( "%s is already a host\n", dsp)
    return
  }

  m := new( HostRecord)
  m.host = h
  m.group = g
  status[dsp] = m
}

func NewCheck( hostdsp string, checkdsp string, ops map[string]string) {
  if !hostExists( hostdsp) {
    log.Printf( "%s does not exist", hostdsp)
    return
  }
  if status[hostdsp].check == nil {
    status[hostdsp].check = make( map[string]*CheckRecord)
  }
  m := new( CheckRecord)
  m.nextRun = time.Now()
  m.freq[0],_ = time.ParseDuration( "90s")
  m.freq[1],_ = time.ParseDuration( "45s")
  m.freq[2],_ = time.ParseDuration( "30s")
  m.freq[3],_ = time.ParseDuration( "30s")
  status[hostdsp].check[checkdsp] = m
}

func PrintHosts() {
  for  k,v := range status {
    log.Printf( "%s %s\n", k, v)
  }
}
